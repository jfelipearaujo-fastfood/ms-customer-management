package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/docker/go-connections/nat"
	"github.com/jfelipearaujo/testcontainers/pkg/container"
	"github.com/jfelipearaujo/testcontainers/pkg/container/localstack"
	"github.com/jfelipearaujo/testcontainers/pkg/container/postgres"
	"github.com/jfelipearaujo/testcontainers/pkg/network"
	"github.com/jfelipearaujo/testcontainers/pkg/state"
	"github.com/jfelipearaujo/testcontainers/pkg/testsuite"
	"github.com/testcontainers/testcontainers-go"
)

type feature struct {
	HostApi    string
	CustomerId string

	StatusCode int
}

var testState = state.NewState[feature]()

var containers = container.NewGroup()

func TestFeatures(t *testing.T) {
	testsuite.NewTestSuite(t,
		initializeScenario,
		testsuite.WithPaths("features"),
		testsuite.WithConcurrency(0),
	)
}

func initializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ntwrkDefinition := network.NewNetwork()

		network, err := ntwrkDefinition.Build(ctx)
		if err != nil {
			return ctx, fmt.Errorf("failed to build the network: %w", err)
		}

		pgDefinition := container.NewContainerDefinition(
			container.WithNetwork(ntwrkDefinition.Alias, network),
			postgres.WithPostgresContainer(),
			container.WithEnvVars(map[string]string{
				"POSTGRES_DB":       "customer_db",
				"POSTGRES_USER":     "customer",
				"POSTGRES_PASSWORD": "customer",
			}),
			container.WithFiles(postgres.BasePath, "./testdata/init-db.sql"),
		)

		pgContainer, err := pgDefinition.BuildContainer(ctx)
		if err != nil {
			return ctx, err
		}

		localStackDefinition := container.NewContainerDefinition(
			localstack.WithLocalStackContainer(),
			container.WithNetwork(ntwrkDefinition.Alias, network),
			container.WithExecutableFiles(localstack.BasePath, "./testdata/init-sm.sh", "./testdata/z-init.sh"),
		)

		localStackContainer, err := localStackDefinition.BuildContainer(ctx)
		if err != nil {
			return ctx, err
		}

		apiDefinition := container.NewContainerDefinition(
			container.WithNetwork(ntwrkDefinition.Alias, network),
			container.WithDockerfile(testcontainers.FromDockerfile{
				Context:    "../",
				Dockerfile: "Dockerfile",
				KeepImage:  false,
			}),
			container.WithEnvVars(map[string]string{
				"API_PORT":              "8080",
				"API_ENV_NAME":          "development",
				"API_VERSION":           "v1",
				"DB_URL":                "todo",
				"DB_URL_SECRET_NAME":    "db-secret-url",
				"AWS_ACCESS_KEY_ID":     "test",
				"AWS_SECRET_ACCESS_KEY": "test",
				"AWS_REGION":            "us-east-1",
				"AWS_BASE_ENDPOINT":     fmt.Sprintf("http://%s:4566", ntwrkDefinition.Alias),
			}),
			container.WithExposedPorts("8080"),
			container.WithWaitingForLog("Server started", 10*time.Second),
		)

		apiContainer, err := apiDefinition.BuildContainer(ctx)
		if err != nil {
			return ctx, err
		}

		host, err := apiContainer.Host(ctx)
		if err != nil {
			return ctx, fmt.Errorf("failed to get the host: %w", err)
		}

		port, err := container.GetMappedPort(ctx, apiContainer, nat.Port("8080/tcp"))
		if err != nil {
			return ctx, err
		}

		containers[sc.Id] = container.BuildGroupContainer(
			container.WithDockerContainer(pgContainer),
			container.WithDockerContainer(apiContainer),
			container.WithDockerContainer(localStackContainer),
		)

		feat := testState.Retrieve(ctx)
		feat.HostApi = fmt.Sprintf("http://%s:%s", host, port)

		return testState.Enrich(ctx, feat), nil
	})

	ctx.Step(`^I have a customer account$`, iHaveACustomerAccount)
	ctx.Step(`^I have a non existent customer account$`, iHaveANonExistentCustomerAccount)
	ctx.Step(`^I delete the customer account$`, iDeleteTheCustomerAccount)
	ctx.Step(`^the customer account should be deleted$`, theCustomerAccountShouldBeDeleted)
	ctx.Step(`^the customer account should not be deleted$`, theCustomerAccountShouldNotBeDeleted)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		group := containers[sc.Id]

		return container.DestroyGroup(ctx, group)
	})
}

func iHaveACustomerAccount(ctx context.Context) (context.Context, error) {
	feat := testState.Retrieve(ctx)
	feat.CustomerId = "19b5408e-8ee2-47d4-953b-196d41f1e367"
	return testState.Enrich(ctx, feat), nil
}

func iHaveANonExistentCustomerAccount(ctx context.Context) (context.Context, error) {
	feat := testState.Retrieve(ctx)
	feat.CustomerId = "e5ccb0dc-f1e9-4ed7-aa8d-49bbac60786e"
	return testState.Enrich(ctx, feat), nil
}

func iDeleteTheCustomerAccount(ctx context.Context) (context.Context, error) {
	feat := testState.Retrieve(ctx)

	client := &http.Client{}
	route := fmt.Sprintf("%s/api/v1/customers/%s", feat.HostApi, feat.CustomerId)

	req, err := http.NewRequest(http.MethodDelete, route, nil)
	if err != nil {
		return ctx, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return ctx, err
	}

	feat.StatusCode = resp.StatusCode

	return testState.Enrich(ctx, feat), nil
}

func theCustomerAccountShouldBeDeleted(ctx context.Context) (context.Context, error) {
	feat := testState.Retrieve(ctx)

	if feat.StatusCode != http.StatusNoContent {
		return ctx, fmt.Errorf("expected status code %d, got %d", http.StatusNoContent, feat.StatusCode)
	}

	return ctx, nil
}

func theCustomerAccountShouldNotBeDeleted(ctx context.Context) (context.Context, error) {
	feat := testState.Retrieve(ctx)

	if feat.StatusCode != http.StatusNotFound {
		return ctx, fmt.Errorf("expected status code %d, got %d", http.StatusNotFound, feat.StatusCode)
	}

	return ctx, nil
}
