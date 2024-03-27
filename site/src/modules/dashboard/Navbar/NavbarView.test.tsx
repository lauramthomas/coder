import { screen } from "@testing-library/react";
import type { ProxyContextValue } from "contexts/ProxyContext";
import { MockPrimaryWorkspaceProxy, MockUser } from "testHelpers/entities";
import { renderWithAuth } from "testHelpers/renderHelpers";
import { Language as navLanguage, NavbarView } from "./NavbarView";

const proxyContextValue: ProxyContextValue = {
  proxy: {
    preferredPathAppURL: "",
    preferredWildcardHostname: "",
    proxy: MockPrimaryWorkspaceProxy,
  },
  isLoading: false,
  isFetched: true,
  setProxy: jest.fn(),
  clearProxy: jest.fn(),
  refetchProxyLatencies: jest.fn(),
  proxyLatencies: {},
};

describe("NavbarView", () => {
  const noop = jest.fn();

  it("workspaces nav link has the correct href", async () => {
    renderWithAuth(
      <NavbarView
        proxyContextValue={proxyContextValue}
        user={MockUser}
        onSignOut={noop}
        canViewAuditLog
        canViewDeployment
        canViewAllUsers
        canViewHealth
        canViewInsights
      />,
    );
    const workspacesLink = await screen.findByText(navLanguage.workspaces);
    expect((workspacesLink as HTMLAnchorElement).href).toContain("/workspaces");
  });

  it("templates nav link has the correct href", async () => {
    renderWithAuth(
      <NavbarView
        proxyContextValue={proxyContextValue}
        user={MockUser}
        onSignOut={noop}
        canViewAuditLog
        canViewDeployment
        canViewAllUsers
        canViewHealth
        canViewInsights
      />,
    );
    const templatesLink = await screen.findByText(navLanguage.templates);
    expect((templatesLink as HTMLAnchorElement).href).toContain("/templates");
  });

  it("users nav link has the correct href", async () => {
    renderWithAuth(
      <NavbarView
        proxyContextValue={proxyContextValue}
        user={MockUser}
        onSignOut={noop}
        canViewAuditLog
        canViewDeployment
        canViewAllUsers
        canViewHealth
        canViewInsights
      />,
    );
    const userLink = await screen.findByText(navLanguage.users);
    expect((userLink as HTMLAnchorElement).href).toContain("/users");
  });

  it("audit nav link has the correct href", async () => {
    renderWithAuth(
      <NavbarView
        proxyContextValue={proxyContextValue}
        user={MockUser}
        onSignOut={noop}
        canViewAuditLog
        canViewDeployment
        canViewAllUsers
        canViewHealth
        canViewInsights
      />,
    );
    const auditLink = await screen.findByText(navLanguage.audit);
    expect((auditLink as HTMLAnchorElement).href).toContain("/audit");
  });

  it("deployment nav link has the correct href", async () => {
    renderWithAuth(
      <NavbarView
        proxyContextValue={proxyContextValue}
        user={MockUser}
        onSignOut={noop}
        canViewAuditLog
        canViewDeployment
        canViewAllUsers
        canViewHealth
        canViewInsights
      />,
    );
    const auditLink = await screen.findByText(navLanguage.deployment);
    expect((auditLink as HTMLAnchorElement).href).toContain(
      "/deployment/general",
    );
  });
});
