import { test, expect } from "@playwright/test";

test("has title", async ({ page }) => {
  await page.goto("");
  await expect(page).toHaveTitle(/Let Him Cook/);
});

test("login", async ({ page }) => {
  await page.goto("");
  await page.locator("#admin-button").click();
  await page.getByPlaceholder("Passwort", { exact: true }).fill("admin");
  await page.getByRole("button", { name: "Anmelden" }).click();
  await expect(page.getByRole("button", { name: "Abmelden" })).toBeVisible();
  await expect(
    page.getByPlaceholder("Passwort", { exact: true }),
  ).toBeDisabled();
  await expect(page.getByRole("button", { name: "Anmelden" })).toBeDisabled();
});
