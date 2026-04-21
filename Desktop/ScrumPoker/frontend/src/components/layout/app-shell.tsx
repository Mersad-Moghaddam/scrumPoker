import { PropsWithChildren } from "react";

import { ThemeToggle } from "../theme-toggle";

export function AppShell({ children }: PropsWithChildren) {
  return (
    <div className="min-h-screen bg-app px-4 py-6 text-app-text sm:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl">
        <div className="mb-6 flex items-center justify-between">
          <div>
            <p className="text-sm text-app-muted">اسکرام پوکر فارسی</p>
            <h1 className="text-2xl font-bold">برآورد تیمی، ساده و هم‌زمان</h1>
          </div>
          <ThemeToggle />
        </div>
        {children}
      </div>
    </div>
  );
}
