import { test, expect, Page } from "@playwright/test";

async function login(page: Page) {
  await page.goto("");
  await page.locator("#admin-button").click();
  await page.getByPlaceholder("Passwort", { exact: true }).fill("admin");
  await page.getByRole("button", { name: "Anmelden" }).click();
  await expect(page.getByRole("button", { name: "Abmelden" })).toBeVisible();
  await expect(
    page.getByPlaceholder("Passwort", { exact: true }),
  ).toBeDisabled();
  await expect(page.getByRole("button", { name: "Anmelden" })).toBeDisabled();
  await page.goto("");
}

async function createRecipe(page: Page) {
  await page.locator("#new-recipe-button").click();

  await page.getByPlaceholder("Titel").fill("Naan");
  await page
    .getByPlaceholder("Beschreibung")
    .fill("Indisches Fladenbrot - zubereitet in der Pfanne");
  await page.getByPlaceholder("Zubereitungszeit (Minuten)").fill("180");

  await page.getByPlaceholder("Zutaten").fill(`**Rezept für 10 Stück**
- 560 g Maida-Mehl
- 10 g Salz
- 5 g Backpulver
- 8 g Streuzucker
- 150 ml Milch
- 135 ml Wasser
- 1 Ei
- 1 EL Pflanzenöl`);

  await page.getByPlaceholder("Anleitung")
    .fill(`1. Das Mehl und das Salz in eine Schüssel sieben und in die Mitte eine Vertiefung drücken
2. Backpulver, Zucker, Milch, Wasser und Ei verrühren
3. Die Mischung dem Mehl langsam untermischen und alles 5 Minuten kneten, anschließend 5 Minuten ruhen lassen
4. Das Öl über den Teig träufeln und diesen solange kneten bis das Öl eingearbeitet ist
5. Den Teig in eine saubere Schüssel geben, mit einem feuchten Küchentuch abdecken und 2 Stunden ruhen lassen
6. Aus dem Teig 70g Kugeln formen, diese dann für 30 Minuten auf einem eingeölten und mit Frischhaltefolie abgedeckten Backblech ruhen lassen
7. Nacheinander die leicht eingeölten Kugeln auf der Arbeitsfläche ausrollen und in einer heißen Pfanne zubereiten`);

  await page.getByRole("button", { name: "Rezept erstellen" }).click();

  await expect(page.getByText("1 Rezept")).toBeVisible();
  await expect(page.locator(".recipe-list-item")).toHaveCount(1);
  await expect(page.getByText("Naan")).toBeVisible();
  await expect(page.getByText("180 Minuten")).toBeVisible();
  await expect(
    page.getByText("Indisches Fladenbrot - zubereitet in der Pfanne"),
  ).toBeVisible();
}

async function deleteRecipe(page: Page) {
  await page.goto("");
  await page.locator("#recipe-delete-button").click();
  await page.getByRole("button", { name: "Löschen bestätigen" }).click();
  await page.goto("");
  await expect(page.getByText("Keine Rezepte")).toBeVisible();
  await expect(page.locator(".recipe-list-item")).toHaveCount(0);
}

test("has title", async ({ page }) => {
  await page.goto("");
  await expect(page).toHaveTitle(/Let Him Cook/);
});

test("login", async ({ page }) => {
  await login(page);
});

test("create recipe", async ({ page }) => {
  await login(page);
  await createRecipe(page);
  await deleteRecipe(page);
});
