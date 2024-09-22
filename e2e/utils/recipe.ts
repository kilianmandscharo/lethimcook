import { Page, expect } from "@playwright/test";
import { clickButtonByTitle, navigateToRecipePage } from "./utils";

export type Recipe = typeof recipe;

export const recipe = {
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

export const editedRecipe = {
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

export async function createRecipe(page: Page, pending?: boolean) {
    await clickButtonByTitle(page, "Neues Rezept");
    await page.getByPlaceholder("Titel").fill(recipe.title);
    await page.getByPlaceholder("Beschreibung").fill(recipe.description);
    await page
        .getByPlaceholder("Zubereitungszeit")
        .fill(recipe.duration);
    await page.getByPlaceholder("Autor").fill(recipe.author);
    await page.getByPlaceholder("Quelle").fill(recipe.source);
    await page.getByPlaceholder("Tags").fill(recipe.tags);
    await page.getByPlaceholder("Zutaten").fill(recipe.ingredients);
    await page.getByPlaceholder("Anleitung").fill(recipe.instructions);
    await page
        .getByRole("button", {
            name: pending ? "Rezept einreichen" : "Rezept erstellen",
        })
        .click();
}

export async function checkRecipeList(page: Page, testRecipe: Recipe) {
    await expect(page.getByText("1 Rezept")).toBeVisible();
    await expect(page.locator(".recipe-list-item")).toHaveCount(1);
    await expect(page.getByText(testRecipe.title)).toBeVisible();
    await expect(page.getByText(`${testRecipe.duration} Minuten`)).toBeVisible();
    await expect(page.getByText(testRecipe.description)).toBeVisible();
}

export async function checkRecipePage(page: Page, testRecipe: Recipe) {
    await expect(page.getByText(testRecipe.title)).toBeVisible();
    await expect(page.getByText(`${testRecipe.duration} Minuten`)).toBeVisible();
    await expect(page.getByText(`${testRecipe.author}`)).toBeVisible();
    await expect(page.getByText(`${testRecipe.source}`)).toBeVisible();
    await expect(page.getByText(testRecipe.description)).toBeVisible();
}


export async function editRecipe(page: Page) {
    await page.getByLabel("Titel").fill(editedRecipe.title);
    await page
        .getByLabel("Zubereitungszeit (Minuten)")
        .fill(editedRecipe.duration);
    await page.getByLabel("Beschreibung").fill(editedRecipe.description);
    await page.getByPlaceholder("Autor").fill(editedRecipe.author);
    await page.getByPlaceholder("Quelle").fill(editedRecipe.source);
    await page.getByRole("button", { name: "Rezept aktualisieren" }).click();
}

export async function deleteRecipe(page: Page) {
    await page.goto("");
    await navigateToRecipePage(page, recipe.title);
    await clickButtonByTitle(page, "Rezept löschen");
    page.on("dialog", (dialog) => dialog.accept());
    await clickButtonByTitle(page, "Rezept löschen");
    await expect(page.getByText("Keine Rezepte")).toBeVisible();
    await expect(page.locator(".recipe-list-item")).toHaveCount(0);
}

