import { fireEvent, render, screen } from "@testing-library/react";

import { ThemeProvider, themeStorageKey, useTheme } from "./use-theme";

function ThemeProbe() {
  const { theme, toggleTheme } = useTheme();
  return (
    <>
      <span>{theme}</span>
      <button onClick={toggleTheme}>toggle</button>
    </>
  );
}

describe("ThemeProvider", () => {
  it("persists the selected theme", () => {
    render(
      <ThemeProvider>
        <ThemeProbe />
      </ThemeProvider>
    );

    expect(screen.getByText("dark")).toBeInTheDocument();

    fireEvent.click(screen.getByText("toggle"));

    expect(screen.getByText("light")).toBeInTheDocument();
    expect(localStorage.getItem(themeStorageKey)).toBe("light");
    expect(document.documentElement.classList.contains("dark")).toBe(false);
  });
});
