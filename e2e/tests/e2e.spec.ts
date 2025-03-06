import { test, expect } from "@playwright/test";
import {
    clickButtonByTitle,
    navigateToEditPage,
    navigateToHomePage,
    navigateToAdminPage,
    navigateToRecipePage,
    navigateToNewRecipePage,
} from "../utils/utils";
import {
    checkRecipeList,
    checkRecipePage,
    createRecipe,
    deleteRecipe,
    editRecipe,
    naan,
    naanEdited,
    butterChicken,
} from "../utils/recipe";
import {
    changePassword,
    login,
    logout,
    testInvalidPassword,
    testInvalidPasswordChange,
    testPasswordChangeTooShort,
} from "../utils/admin";

test("check title", async ({ page }) => {
    await page.goto("");
    await expect(page).toHaveTitle(/Let Him Cook/);
});

test("admin tests", async ({ page }) => {
    await page.goto("");
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
});

// test("save wip recipe in local storage", async ({ page }) => {
//     await page.goto("");
//     await navigateToAdminPage(page);
//     await login(page, "admin");
//     await navigateToHomePage(page);
//     await navigateToNewRecipePage(page);
//     await fillInputByPlaceholder(page, "Zutaten", "Test Zutaten");
//     await page.keyboard.press("Tab");
//     await navigateToHomePage(page);
//     await navigateToNewRecipePage(page);
//     await expect(async () => {
//         const value = await page.getByPlaceholder("Zutaten").inputValue();
//         console.log("VALUE:", value);
//         expect(page.getByPlaceholder("Zutaten")).toHaveValue("Test Zutaten");
//     }).toPass({ timeout: 5_000, intervals: [1_000] });
// });
//

test("create recipe", async ({ page }) => {
    await page.goto("");
    await navigateToAdminPage(page);
    await login(page, "admin");
    await navigateToHomePage(page);
    await createRecipe(page, naan);
    await checkRecipeList(page, naan);
    await navigateToRecipePage(page, naan.title);
    await checkRecipePage(page, naan);
});

test("create second recipe", async ({ page }) => {
    await page.goto("");
    await navigateToAdminPage(page);
    await login(page, "admin");
    await navigateToHomePage(page);
    await createRecipe(page, butterChicken);
    await checkRecipeList(page, butterChicken);
    await navigateToRecipePage(page, butterChicken.title);
    await checkRecipePage(page, butterChicken);
});

test("select dialog", async ({ page }) => {
    await page.goto("");
    await navigateToNewRecipePage(page);
    const input = page.getByPlaceholder("Zutaten");
    await expect(input).toBeVisible();
    await input.click();
    await expect(input).toBeFocused();

    await input.pressSequentially("[!indisch]");

    await expect(page.locator("#select-dialog")).toBeVisible();
    await expect(page.getByText("Rezept auswählen")).toBeVisible();
    await expect(
        page.locator("#options-container").getByText("Butter Chicken"),
    ).toBeVisible();
    await expect(
        page.locator("#options-container").getByText("Naan"),
    ).toBeVisible();
    await page
        .locator("#options-container")
        .getByText("Butter Chicken")
        .click();
    await expect(page.getByPlaceholder("Zutaten")).toHaveValue(
        /\[Butter Chicken\]\(http:\/\/127\.0\.0\.1:8080\/recipe\/\d+\)/,
    );
});

test("recipe link completion", async ({ page }) => {
    await page.goto("");
    await navigateToNewRecipePage(page);
    const input = page.getByPlaceholder("Zutaten");
    await expect(input).toBeVisible();
    await input.click();
    await expect(input).toBeFocused();

    await input.pressSequentially("[!Naan]");

    await expect(page.getByPlaceholder("Zutaten")).toHaveValue(
        /\[Naan\]\(http:\/\/127\.0\.0\.1:8080\/recipe\/\d+\)/,
    );
});

test("edit recipe", async ({ page }) => {
    await page.goto("");
    await navigateToAdminPage(page);
    await login(page, "admin");
    await navigateToHomePage(page);
    await navigateToRecipePage(page, naan.title);
    await navigateToEditPage(page);
    await editRecipe(page, naanEdited);
    await checkRecipePage(page, naanEdited);
    await navigateToHomePage(page);
    await checkRecipeList(page, naanEdited);
});

test("delete recipes", async ({ page }) => {
    await page.goto("");
    await navigateToAdminPage(page);
    await login(page, "admin");
    await deleteRecipe(page, naan);
    await deleteRecipe(page, butterChicken);
    await navigateToHomePage(page);
    await expect(page.getByText("Keine Rezepte")).toBeVisible();
    await expect(page.locator(".recipe-list-item")).toHaveCount(0);
});

test("create pending recipe", async ({ page }) => {
    page.on("dialog", (dialog) => dialog.accept());

    await page.goto("");
    await createRecipe(page, naan, true);

    navigateToHomePage(page);
    await expect(page.getByText("Keine Rezepte")).toBeVisible();
    await expect(page.locator(".recipe-list-item")).toHaveCount(0);

    await navigateToAdminPage(page);
    await login(page, "admin");

    navigateToHomePage(page);
    await checkRecipeList(page, naan);

    await navigateToRecipePage(page, naan.title);
    await clickButtonByTitle(page, "Rezept akzeptieren");

    await navigateToAdminPage(page);
    await logout(page);

    navigateToHomePage(page);
    await checkRecipeList(page, naan);

    await navigateToAdminPage(page);
    await login(page, "admin");
    navigateToHomePage(page);
    await navigateToRecipePage(page, naan.title);
    const button = page.getByTitle("Rezept auf 'ausstehend' setzen");
    await expect(button).toBeVisible();
    await button.click();

    await navigateToAdminPage(page);
    await logout(page);

    navigateToHomePage(page);
    await expect(page.getByText("Keine Rezepte")).toBeVisible();
    await expect(page.locator(".recipe-list-item")).toHaveCount(0);
});

test("clean up", async ({ page }) => {
    await page.goto("");
    await navigateToAdminPage(page);
    await login(page, "admin");
    await navigateToHomePage(page);
    await navigateToRecipePage(page, naan.title);
    await clickButtonByTitle(page, "Rezept ablehnen");
    page.on("dialog", (dialog) => dialog.accept());
    await clickButtonByTitle(page, "Rezept ablehnen");
});
