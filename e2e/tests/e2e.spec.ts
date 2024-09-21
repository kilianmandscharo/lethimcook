import { test, expect } from "@playwright/test";
import { clickButtonByTitle, navigateToEditPage, navigateToHomePage, navigateToAdminPage, navigateToRecipePage } from "../utils/utils";
import { checkRecipeList, checkRecipePage, createRecipe, deleteRecipe, editRecipe, editedRecipe, recipe } from "../utils/recipe";
import { changePassword, login, logout, testInvalidPassword, testInvalidPasswordChange, testPasswordChangeTooShort } from "../utils/admin";

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

test("create recipe", async ({ page }) => {
    await page.goto("");
    await navigateToAdminPage(page);
    await login(page, "admin");
    await navigateToHomePage(page);
    await createRecipe(page);
    await checkRecipeList(page, recipe);
    await navigateToRecipePage(page, recipe.title);
    await checkRecipePage(page, recipe);
});

test("edit recipe", async ({ page }) => {
    await page.goto("");
    await navigateToAdminPage(page);
    await login(page, "admin");
    await navigateToHomePage(page);
    await navigateToRecipePage(page, recipe.title);
    await navigateToEditPage(page);
    await editRecipe(page);
    await checkRecipePage(page, editedRecipe);
    await navigateToHomePage(page);
    await checkRecipeList(page, editedRecipe);
});

test("delete recipe", async ({ page }) => {
    await page.goto("");
    await navigateToAdminPage(page);
    await login(page, "admin");
    await deleteRecipe(page);
});

test("create pending recipe", async ({ page }) => {
    page.on("dialog", (dialog) => dialog.accept());

    await page.goto("");
    await createRecipe(page, true);

    navigateToHomePage(page);
    await expect(page.getByText("Keine Rezepte")).toBeVisible();
    await expect(page.locator(".recipe-list-item")).toHaveCount(0);

    await navigateToAdminPage(page);
    await login(page, "admin");

    navigateToHomePage(page);
    await checkRecipeList(page, recipe);

    await navigateToRecipePage(page, recipe.title);
    await clickButtonByTitle(page, "Rezept akzeptieren");

    await navigateToAdminPage(page);
    await logout(page);

    navigateToHomePage(page);
    await checkRecipeList(page, recipe);

    await navigateToAdminPage(page);
    await login(page, "admin");
    navigateToHomePage(page);
    await navigateToRecipePage(page, recipe.title);
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
    await navigateToRecipePage(page, recipe.title);
    await clickButtonByTitle(page, "Rezept ablehnen");
    page.on("dialog", (dialog) => dialog.accept());
    await clickButtonByTitle(page, "Rezept ablehnen");
});
