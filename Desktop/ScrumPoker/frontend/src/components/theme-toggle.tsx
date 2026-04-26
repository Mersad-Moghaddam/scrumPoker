import { MoonStar, SunMedium } from "lucide-react";

import { useTheme } from "../hooks/use-theme";
import { Button } from "./ui/button";

export function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <Button
      variant="secondary"
      onClick={toggleTheme}
      className="px-3"
      aria-label={theme === "dark" ? "حالت روشن" : "حالت تیره"}
    >
      {theme === "dark" ? <SunMedium size={16} /> : <MoonStar size={16} />}
    </Button>
  );
}
