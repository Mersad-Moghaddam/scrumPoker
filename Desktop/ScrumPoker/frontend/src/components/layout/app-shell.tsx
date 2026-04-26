import { PropsWithChildren } from "react";

import { ThemeToggle } from "../theme-toggle";

export function AppShell({ children }: PropsWithChildren) {
  return (
    <div className="min-h-screen bg-app px-4 py-4 text-app-text sm:px-6">
      <div className="mx-auto max-w-5xl">
        <div className="mb-4 flex items-center justify-between">
          <h1 className="text-lg font-semibold">اسکرام پوکر</h1>
          <ThemeToggle />
        </div>
        {children}
      </div>
    </div>
  );
}
