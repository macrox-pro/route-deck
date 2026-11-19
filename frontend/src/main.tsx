import "@/styles/globals.css";

import React from "react";
import ReactDOM from "react-dom/client";
import { HashRouter, Navigate, Route, Routes } from "react-router-dom";
import { AppShell } from "@/components/application/app-shell/app-shell";
import { RouteProvider } from "@/providers/route-provider";
import { ThemeProvider } from "@/providers/theme-provider";
import { ActivityPage } from "@/pages/activity";
import { ApplicationsPage } from "@/pages/applications";
import { DashboardPage } from "@/pages/dashboard";
import { AboutDialog } from "@/pages/dialogs/about-dialog";
import { GuidesDialog } from "@/pages/dialogs/guides-dialog";
import { HelpersDialog } from "@/pages/dialogs/helpers-dialog";
import { ProfileEditorDialog } from "@/pages/dialogs/profile-editor-dialog";
import { ProxyEditorDialog } from "@/pages/dialogs/proxy-editor-dialog";
import { ProfilesPage } from "@/pages/profiles";
import { ProxiesPage } from "@/pages/proxies";
import { SettingsPage } from "@/pages/settings";

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <HashRouter>
            <RouteProvider>
                <ThemeProvider>
                    <Routes>
                        <Route path="dialog/profile" element={<ProfileEditorDialog />} />
                        <Route path="dialog/profile/:id" element={<ProfileEditorDialog />} />
                        <Route path="dialog/proxy" element={<ProxyEditorDialog />} />
                        <Route path="dialog/proxy/:id" element={<ProxyEditorDialog />} />
                        <Route path="dialog/helpers/:applicationId" element={<HelpersDialog />} />
                        <Route path="dialog/about" element={<AboutDialog />} />
                        <Route path="dialog/guides" element={<GuidesDialog />} />
                        <Route path="dialog/guides/:step" element={<GuidesDialog />} />
                        <Route element={<AppShell />}>
                            <Route index element={<DashboardPage />} />
                            <Route path="applications" element={<ApplicationsPage />} />
                            <Route path="proxies" element={<ProxiesPage />} />
                            <Route path="profiles" element={<ProfilesPage />} />
                            <Route path="activity" element={<ActivityPage />} />
                            <Route path="settings" element={<SettingsPage />} />
                            <Route path="*" element={<Navigate to="/" replace />} />
                        </Route>
                    </Routes>
                </ThemeProvider>
            </RouteProvider>
        </HashRouter>
    </React.StrictMode>,
);
