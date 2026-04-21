import { MoonStar, SunMedium } from "lucide-react";

import { useTheme } from "../hooks/use-theme";
import { Button } from "./ui/button";

export function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <Button variant="secondary" onClick={toggleTheme} className="gap-2">
      {theme === "dark" ? <SunMedium size={16} /> : <MoonStar size={16} />}
      <span>{theme === "dark" ? "حالت روشن" : "حالت تیره"}</span>
    </Button>
  );
}
