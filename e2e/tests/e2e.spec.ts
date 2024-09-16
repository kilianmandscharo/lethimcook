import { test, expect, Page } from "@playwright/test";

const recipe = {
  title: "Naan",
  description: "Indisches Fladenbrot - zubereitet in der Pfanne",
  duration: "180",
  author: "Phillip Jeffries",
  source: "Indische Küche Dishoom",
  tags: "indisch, Beilage",
  ingredients: `**Rezept für 10 Stück**
- 560 g Maida-Mehl
- 10 g Salz
- 5 g Backpulver
- 8 g Streuzucker
- 150 ml Milch
- 135 ml Wasser
- 1 Ei
- 1 EL Pflanzenöl`,
  instructions: `1. Das Mehl und das Salz in eine Schüssel sieben und in die Mitte eine Vertiefung drücken
2. Backpulver, Zucker, Milch, Wasser und Ei verrühren
3. Die Mischung dem Mehl langsam untermischen und alles 5 Minuten kneten, anschließend 5 Minuten ruhen lassen
4. Das Öl über den Teig träufeln und diesen solange kneten bis das Öl eingearbeitet ist
5. Den Teig in eine saubere Schüssel geben, mit einem feuchten Küchentuch abdecken und 2 Stunden ruhen lassen
6. Aus dem Teig 70g Kugeln formen, diese dann für 30 Minuten auf einem eingeölten und mit Frischhaltefolie abgedeckten Backblech ruhen lassen
7. Nacheinander die leicht eingeölten Kugeln auf der Arbeitsfläche ausrollen und in einer heißen Pfanne zubereiten`,
};

const editedRecipe = {
  title: "Naanbrot",
  description: "Leckeres indisches Fladenbrot - zubereitet in der Pfanne",
  duration: "210",
  author: "Dale Cooper",
  source: "Indische Küche Dishoom (Buch)",
  tags: "Indien",
  ingredients: `**Rezept für 20 Stück**
- 1120 g Maida-Mehl
- 20 g Salz
- 10 g Backpulver
- 16 g Streuzucker
- 300 ml Milch
- 270 ml Wasser
- 2 Eier
- 2 EL Pflanzenöl`,
  instructions: `1. Das Mehl und das Salz in eine Schüssel sieben und in die Mitte eine Vertiefung drücken
2. Backpulver, Zucker, Milch, Wasser und Ei verrühren
3. Die Mischung dem Mehl langsam untermischen und alles 5 Minuten kneten, anschließend 5 Minuten ruhen lassen
4. Das Öl über den Teig träufeln und diesen solange kneten bis das Öl eingearbeitet ist
5. Den Teig in eine saubere Schüssel geben, mit einem feuchten Küchentuch abdecken und 2 Stunden ruhen lassen
6. Aus dem Teig 70g Kugeln formen, diese dann für 30 Minuten auf einem eingeölten und mit Frischhaltefolie abgedeckten Backblech ruhen lassen
7. Nacheinander die leicht eingeölten Kugeln auf der Arbeitsfläche ausrollen und in einer heißen Pfanne zubereiten`,
};

type Recipe = typeof recipe;

async function navigateToAdminPage(page: Page) {
  await page.goto("");
  await page.locator("#admin-button").click();
}

async function login(page: Page, password: string) {
  await page.getByPlaceholder("Passwort", { exact: true }).fill(password);
  await page.getByRole("button", { name: "Anmelden" }).click();
  await expect(page.getByRole("button", { name: "Abmelden" })).toBeVisible();
  await expect(
    page.getByPlaceholder("Passwort", { exact: true }),
  ).toBeDisabled();
  await expect(page.getByRole("button", { name: "Anmelden" })).toBeDisabled();
}

async function logout(page: Page) {
  await page.getByRole("button", { name: "Abmelden" }).click();
  await expect(
    page.getByRole("button", { name: "Abmelden" }),
  ).not.toBeVisible();
  await expect(
    page.getByPlaceholder("Passwort", { exact: true }),
  ).toBeEnabled();
  await expect(page.getByRole("button", { name: "Anmelden" })).toBeEnabled();
}

async function changePassword(
  page: Page,
  oldPassword: string,
  newPassword: string,
) {
  await page.getByPlaceholder("Altes Passwort").fill(oldPassword);
  await page.getByPlaceholder("Neues Passwort").fill(newPassword);
  await page.getByRole("button", { name: "Bestätigen" }).click();
  await expect(page.getByPlaceholder("Altes Passwort")).toBeEmpty();
  await expect(page.getByPlaceholder("Neues Passwort")).toBeEmpty();
}

async function testInvalidPasswordChange(page: Page) {
  await changePassword(page, "invalid", "nimda");
  await expect(
    page.locator("#content").getByText("Falsches Passwort"),
  ).toBeVisible();
}

async function testPasswordChangeTooShort(page: Page) {
  await changePassword(page, "admin", "aaa");
  await expect(
    page.locator("#content").getByText("Minimale Passwortlänge: 5"),
  ).toBeVisible();
}

async function testInvalidPassword(page: Page, invalidPassword: string) {
  await page
    .getByPlaceholder("Passwort", { exact: true })
    .fill(invalidPassword);
  await page.getByRole("button", { name: "Anmelden" }).click();
  await page.waitForLoadState("networkidle");
  await expect(
    page.getByRole("button", { name: "Abmelden" }),
  ).not.toBeVisible();
  await expect(
    page.locator("#content").getByText("Falsches Passwort"),
  ).toBeVisible();
}

async function navigateToHomePage(page: Page) {
  await page.locator("#home-button").click();
}

async function createRecipe(page: Page) {
  await page.locator("#new-recipe-button").click();

  await page.getByPlaceholder("Titel").fill(recipe.title);
  await page.getByPlaceholder("Beschreibung").fill(recipe.description);
  await page
    .getByPlaceholder("Zubereitungszeit (Minuten)")
    .fill(recipe.duration);
  await page.getByPlaceholder("Autor").fill(recipe.author);
  await page.getByPlaceholder("Quelle").fill(recipe.source);
  await page.getByPlaceholder("Tags").fill(recipe.tags);
  await page.getByPlaceholder("Zutaten").fill(recipe.ingredients);
  await page.getByPlaceholder("Anleitung").fill(recipe.instructions);

  await page.getByRole("button", { name: "Rezept erstellen" }).click();
}

async function checkRecipeList(page: Page, testRecipe: Recipe) {
  await expect(page.getByText("1 Rezept")).toBeVisible();
  await expect(page.locator(".recipe-list-item")).toHaveCount(1);
  await expect(page.getByText(testRecipe.title)).toBeVisible();
  await expect(page.getByText(`${testRecipe.duration} Minuten`)).toBeVisible();
  await expect(page.getByText(testRecipe.description)).toBeVisible();
}

async function navigateToRecipePage(page: Page) {
  await page.getByText(recipe.title).click();
}

async function checkRecipePage(page: Page, testRecipe: Recipe) {
  await expect(page.getByText(testRecipe.title)).toBeVisible();
  await expect(
    page.getByText(`Zubereitungszeit: ${testRecipe.duration} Minuten`),
  ).toBeVisible();
  await expect(page.getByText(`Autor: ${testRecipe.author}`)).toBeVisible();
  await expect(page.getByText(`Quelle: ${testRecipe.source}`)).toBeVisible();
  await expect(page.getByText(testRecipe.description)).toBeVisible();
}

async function navigateToEditPage(page: Page) {
  await page.locator("#edit-recipe-button").click();
}

async function editRecipe(page: Page) {
  await page.getByLabel("Titel").fill(editedRecipe.title);
  await page
    .getByLabel("Zubereitungszeit (Minuten)")
    .fill(editedRecipe.duration);
  await page.getByLabel("Beschreibung").fill(editedRecipe.description);
  await page.getByPlaceholder("Autor").fill(editedRecipe.author);
  await page.getByPlaceholder("Quelle").fill(editedRecipe.source);
  await page.getByRole("button", { name: "Rezept aktualisieren" }).click();
}

async function deleteRecipe(page: Page) {
  await page.goto("");
  await page.locator("#recipe-delete-button").click();
  await page.getByRole("button", { name: "Löschen bestätigen" }).click();
  await page.goto("");
  await expect(page.getByText("Keine Rezepte")).toBeVisible();
  await expect(page.locator(".recipe-list-item")).toHaveCount(0);
}

test("e2e test", async ({ page }) => {
  await page.goto("");
  await expect(page).toHaveTitle(/Let Him Cook/);

  await navigateToAdminPage(page);

  await testInvalidPasswordChange(page);
  await testPasswordChangeTooShort(page);

  await testInvalidPassword(page, "invalid");
  await login(page, "admin");
  await logout(page);
  await changePassword(page, "admin", "nimda");
  await testInvalidPassword(page, "admin");
  await login(page, "nimda");
  await logout(page);
  await changePassword(page, "nimda", "admin");
  await testInvalidPassword(page, "nimda");
  await login(page, "admin");

  await navigateToHomePage(page);
  await createRecipe(page);
  await checkRecipeList(page, recipe);
  await navigateToRecipePage(page);
  await checkRecipePage(page, recipe);

  await navigateToEditPage(page);
  await editRecipe(page);
  await checkRecipePage(page, editedRecipe);
  await navigateToHomePage(page);
  await checkRecipeList(page, editedRecipe);

  await deleteRecipe(page);
});
