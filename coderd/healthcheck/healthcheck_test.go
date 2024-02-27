package healthcheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coder/coder/v2/coderd/healthcheck"
	"github.com/coder/coder/v2/coderd/healthcheck/derphealth"
	"github.com/coder/coder/v2/coderd/healthcheck/health"
	"github.com/coder/coder/v2/codersdk"
)

type testChecker struct {
	DERPReport               codersdk.DERPHealthReport
	AccessURLReport          codersdk.AccessURLReport
	WebsocketReport          codersdk.WebsocketReport
	DatabaseReport           codersdk.DatabaseReport
	WorkspaceProxyReport     codersdk.WorkspaceProxyReport
	ProvisionerDaemonsReport codersdk.ProvisionerDaemonsReport
}

func (c *testChecker) DERP(context.Context, *derphealth.ReportOptions) codersdk.DERPHealthReport {
	return c.DERPReport
}

func (c *testChecker) AccessURL(context.Context, *healthcheck.AccessURLReportOptions) codersdk.AccessURLReport {
	return c.AccessURLReport
}

func (c *testChecker) Websocket(context.Context, *healthcheck.WebsocketReportOptions) codersdk.WebsocketReport {
	return c.WebsocketReport
}

func (c *testChecker) Database(context.Context, *healthcheck.DatabaseReportOptions) codersdk.DatabaseReport {
	return c.DatabaseReport
}

func (c *testChecker) WorkspaceProxy(context.Context, *healthcheck.WorkspaceProxyReportOptions) codersdk.WorkspaceProxyReport {
	return c.WorkspaceProxyReport
}

func (c *testChecker) ProvisionerDaemons(context.Context, *healthcheck.ProvisionerDaemonsReportDeps) codersdk.ProvisionerDaemonsReport {
	return c.ProvisionerDaemonsReport
}

func TestHealthcheck(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name            string
		checker         *testChecker
		healthy         bool
		severity        health.Severity
		failingSections []codersdk.HealthSection
	}{{
		name: "OK",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityOK,
			},
		},
		healthy:         true,
		severity:        health.SeverityOK,
		failingSections: []codersdk.HealthSection{},
	}, {
		name: "DERPFail",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityOK,
			},
		},
		healthy:         false,
		severity:        health.SeverityError,
		failingSections: []codersdk.HealthSection{codersdk.HealthSectionDERP},
	}, {
		name: "DERPWarning",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Warnings: []health.Message{{Message: "foobar", Code: "EFOOBAR"}},
				Severity: health.SeverityWarning,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityOK,
			},
		},
		healthy:         true,
		severity:        health.SeverityWarning,
		failingSections: []codersdk.HealthSection{},
	}, {
		name: "AccessURLFail",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  false,
				Severity: health.SeverityWarning,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityOK,
			},
		},
		healthy:         false,
		severity:        health.SeverityWarning,
		failingSections: []codersdk.HealthSection{codersdk.HealthSectionAccessURL},
	}, {
		name: "WebsocketFail",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityOK,
			},
		},
		healthy:         false,
		severity:        health.SeverityError,
		failingSections: []codersdk.HealthSection{codersdk.HealthSectionWebsocket},
	}, {
		name: "DatabaseFail",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityOK,
			},
		},
		healthy:         false,
		severity:        health.SeverityError,
		failingSections: []codersdk.HealthSection{codersdk.HealthSectionDatabase},
	}, {
		name: "ProxyFail",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityOK,
			},
		},
		severity:        health.SeverityError,
		healthy:         false,
		failingSections: []codersdk.HealthSection{codersdk.HealthSectionWorkspaceProxy},
	}, {
		name: "ProxyWarn",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Warnings: []health.Message{{Message: "foobar", Code: "EFOOBAR"}},
				Severity: health.SeverityWarning,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityOK,
			},
		},
		severity:        health.SeverityWarning,
		healthy:         true,
		failingSections: []codersdk.HealthSection{},
	}, {
		name: "ProvisionerDaemonsFail",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityError,
			},
		},
		severity:        health.SeverityError,
		healthy:         false,
		failingSections: []codersdk.HealthSection{codersdk.HealthSectionProvisionerDaemons},
	}, {
		name: "ProvisionerDaemonsWarn",
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  true,
				Severity: health.SeverityOK,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityWarning,
				Warnings: []health.Message{{Message: "foobar", Code: "EFOOBAR"}},
			},
		},
		severity:        health.SeverityWarning,
		healthy:         true,
		failingSections: []codersdk.HealthSection{},
	}, {
		name:    "AllFail",
		healthy: false,
		checker: &testChecker{
			DERPReport: codersdk.DERPHealthReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			AccessURLReport: codersdk.AccessURLReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			WebsocketReport: codersdk.WebsocketReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			DatabaseReport: codersdk.DatabaseReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			WorkspaceProxyReport: codersdk.WorkspaceProxyReport{
				Healthy:  false,
				Severity: health.SeverityError,
			},
			ProvisionerDaemonsReport: codersdk.ProvisionerDaemonsReport{
				Severity: health.SeverityError,
			},
		},
		severity: health.SeverityError,
		failingSections: []codersdk.HealthSection{
			codersdk.HealthSectionDERP,
			codersdk.HealthSectionAccessURL,
			codersdk.HealthSectionWebsocket,
			codersdk.HealthSectionDatabase,
			codersdk.HealthSectionWorkspaceProxy,
			codersdk.HealthSectionProvisionerDaemons,
		},
	}} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			report := healthcheck.Run(context.Background(), &healthcheck.ReportOptions{
				Checker: c.checker,
			})

			assert.Equal(t, c.healthy, report.Healthy)
			assert.Equal(t, c.severity, report.Severity)
			assert.Equal(t, c.failingSections, report.FailingSections)
			assert.Equal(t, c.checker.DERPReport.Healthy, report.DERP.Healthy)
			assert.Equal(t, c.checker.DERPReport.Severity, report.DERP.Severity)
			assert.Equal(t, c.checker.DERPReport.Warnings, report.DERP.Warnings)
			assert.Equal(t, c.checker.AccessURLReport.Healthy, report.AccessURL.Healthy)
			assert.Equal(t, c.checker.AccessURLReport.Severity, report.AccessURL.Severity)
			assert.Equal(t, c.checker.WebsocketReport.Healthy, report.Websocket.Healthy)
			assert.Equal(t, c.checker.WorkspaceProxyReport.Healthy, report.WorkspaceProxy.Healthy)
			assert.Equal(t, c.checker.WorkspaceProxyReport.Warnings, report.WorkspaceProxy.Warnings)
			assert.Equal(t, c.checker.WebsocketReport.Severity, report.Websocket.Severity)
			assert.Equal(t, c.checker.DatabaseReport.Healthy, report.Database.Healthy)
			assert.Equal(t, c.checker.DatabaseReport.Severity, report.Database.Severity)
			assert.NotZero(t, report.Time)
			assert.NotZero(t, report.CoderVersion)
		})
	}
}
