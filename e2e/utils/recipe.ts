import { Dialog, Page, expect } from "@playwright/test";
import { clickButtonByTitle, navigateToRecipePage } from "./utils";

export type Recipe = {
    title: string;
    description: string;
    cookingDuration: string;
    totalDuration: string;
    author: string;
    source: string;
    tags: string;
    ingredients: string;
    instructions: string;
};

export const naan: Recipe = {
    title: "Naan",
    description: "Indisches Fladenbrot - zubereitet in der Pfanne",
    cookingDuration: "100",
    totalDuration: "180",
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

export const naanEdited: Recipe = {
    title: "Naanbrot",
    description: "Leckeres indisches Fladenbrot - zubereitet in der Pfanne",
    cookingDuration: "130",
    totalDuration: "210",
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

export const butterChicken: Recipe = {
    title: "Butter Chicken",
    description: "Der Kultklassiker aus Indien",
    cookingDuration: "80",
    totalDuration: "100",
    author: "Dale Cooper",
    source: "Indische Küche Dishoom (Buch)",
    tags: "indisch, Hähnchen",
    ingredients: `**Für die Soße**
- 35 g Knoblauch
- 175 ML Pflanzenöl
- 20 g Ingwer
- 800 g gehackte Tomaten
- 2 Lorbeerblätter
- 6 gründe Kardamomkapseln
- 2 schwarze Kardamomkapseln
- 2 Zimtstangen
- 2 TL Salz
- 1 1/2 Deggi-Mirch-Chilipulver
- 30 g Butter
- 1 TL Garam Masala
- 20 g Zucker
- 1 EL Honig
- 1 TL Kreuzkümmel
- 1 TL getrockneten Bockshornklee, zerkrümelt
- 1/2 TL frischen Dill

**Für die Marinade**
- 10 g Ingwer
- 20 g gehackten Knoblauch
- 5 g Salz
- 1 TL Deggi-Mirch-Chilipulver
- 1 1/2 TL Kreuzkümmel
- 1/2 TL Garam Masala
- 2 TL Limettensaft
- 2 TL Pflanzenöl
- 75 g griechischen Joghurt

**Restliche Zutaten**
- 700 g Hähnchenfleisch (Brust oder Oberschenkel, ohne Knochen)
- 20 g geschmolzene Butter
- 130 ml Sahne
- Petersilie
- Ingwerstifte`,
    instructions: `**Die Soße (Für optimalen Geschmack am Tag zuvor vorbereiten)**
1. 15 g Knoblauch fein hacken und in einem Topf mit dem Öl goldbraun braten, anschließend den Knoblauch extrahieren und auf einem Küchenpapier abtropfen lassen
2. Den restlichen Knoblauch und Ingwer zu einer Paste zerkleinern und die Tomaten pürieren
3. Bei mittelhoher Hitze im gleichen Öl die Lorbeerblätter, Kardamomkapseln und Zimtstangen eine Minute unter ständigem Rühren anbraten
4. Die Temperatur reduzieren und die Knoblauch-Ingwerpaste etwa 5 Minuten vorsichtig bräunen
5. Nun die Tomaten zusammen mit dem Salz und Chilipulver in den Topf geben, das ganze aufkochen und dann bei niedriger Hitze um die Hälfte einkochen, etwa 30 Minuten
6. Die Zutaten von Schritt 3 aus der Soße fischen, die Butter unterrühren und alles weitere 5 Minuten köcheln
7. Garam Masala, Zucker, Honig, Kreuzkümmel, knusprigen Knoblauch, Bockshornklee und Dill unterrühren und alles weitere 15 Minuten köcheln

**Das Fleisch**
1. Das Fleisch in gleich große Stücke schneiden, alle Zutaten der Marinade pürieren und dann das Fleisch und die Marinade in einer Schüssel vermischen (das Ganze wenn möglich 6-24 Stunden im Kühlschrank ziehen lassen)
2. Das Fleisch auf einem Blech verteilen, mit der geschmolzenen Butter bestreichen und im Ofen (Grillfunktion) bräunen und durchgaren

**Showtime**
1. In einem Topf die Soße aufwärmen, die Sahne und anschließend das Fleisch hineingeben und alles 10 Minuten köcheln
2. Das Ganze mit der Petersilie und den Ingwerstiften zu Reis und [Naan](https://let-him-cook.de/recipe/1) servieren`,
};

export async function createRecipe(
    page: Page,
    recipe: Recipe,
    pending?: boolean,
) {
    await clickButtonByTitle(page, "Neues Rezept");
    await page.getByPlaceholder("Titel").fill(recipe.title);
    await page.getByPlaceholder("Beschreibung").fill(recipe.description);
    await page.getByPlaceholder("Kochzeit").fill(recipe.cookingDuration);
    await page.getByPlaceholder("Gesamtzeit").fill(recipe.totalDuration);
    await page.getByPlaceholder("Autor").fill(recipe.author);
    await page.getByPlaceholder("Quelle").fill(recipe.source);
    await page.getByPlaceholder("Tags").fill(recipe.tags);
    await page.getByPlaceholder("Zutaten").fill(recipe.ingredients);
    await page.getByPlaceholder("Anleitung").fill(recipe.instructions);

    const previewButtonIngredients = page.getByTitle("Vorschau").first();
    await expect(previewButtonIngredients).toBeVisible();
    await previewButtonIngredients.click();

    await expect(page.locator("#preview-modal")).toBeVisible();
    await expect(
        page.locator(".recipe").getByText("Zutaten").first(),
    ).toBeVisible();
    page.keyboard.down("Escape");
    await expect(page.locator("#preview-modal")).not.toBeVisible();

    const previewButtonInstructions = page.getByTitle("Vorschau").nth(1);
    await expect(previewButtonInstructions).toBeVisible();
    await previewButtonInstructions.click();

    await expect(page.locator("#preview-modal")).toBeVisible();
    await expect(page.locator(".recipe").getByText("Anleitung")).toBeVisible();
    page.keyboard.down("Escape");
    await expect(page.locator("#preview-modal")).not.toBeVisible();

    await page
        .getByRole("button", {
            name: pending ? "Rezept einreichen" : "Rezept erstellen",
        })
        .click();
}

export async function checkRecipeList(page: Page, testRecipe: Recipe) {
    await expect(page.getByText(testRecipe.title)).toBeVisible();
    await expect(
        page.getByText(`${testRecipe.totalDuration} Min`),
    ).toBeVisible();
    await expect(page.getByText(testRecipe.description)).toBeVisible();
}

export async function checkRecipePage(page: Page, testRecipe: Recipe) {
    await expect(page.getByText(testRecipe.title)).toBeVisible();
    await expect(
        page.getByText(`${testRecipe.cookingDuration} Minuten`),
    ).toBeVisible();
    await expect(
        page.getByText(`${testRecipe.totalDuration} Minuten`),
    ).toBeVisible();
    await expect(page.getByText(`${testRecipe.author}`)).toBeVisible();
    await expect(page.getByText(`${testRecipe.source}`)).toBeVisible();
    await expect(page.getByText(testRecipe.description)).toBeVisible();
}

export async function editRecipe(page: Page, editedRecipe: Recipe) {
    await page.getByLabel("Titel").fill(editedRecipe.title);
    await page
        .getByLabel("Kochzeit (Minuten)")
        .fill(editedRecipe.cookingDuration);
    await page
        .getByLabel("Gesamtzeit (Minuten)")
        .fill(editedRecipe.totalDuration);
    await page.getByLabel("Beschreibung").fill(editedRecipe.description);
    await page.getByPlaceholder("Autor").fill(editedRecipe.author);
    await page.getByPlaceholder("Quelle").fill(editedRecipe.source);
    await page.getByRole("button", { name: "Rezept aktualisieren" }).click();
}

export async function deleteRecipe(page: Page, recipe: Recipe) {
    await page.goto("");
    await navigateToRecipePage(page, recipe.title);
    await clickButtonByTitle(page, "Rezept löschen");
    const handler = (dialog: Dialog) => dialog.accept();
    page.on("dialog", handler);
    await clickButtonByTitle(page, "Rezept löschen");
    page.off("dialog", handler);
}
