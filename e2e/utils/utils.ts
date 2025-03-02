import { Page, expect } from "@playwright/test";

export async function clickButtonByTitle(page: Page, title: string) {
    const button = page.getByTitle(title);
    await expect(button).toBeVisible();
    await button.click();
}

export async function navigateToRecipePage(page: Page, recipeTitle: string) {
    await page.getByText(recipeTitle).click();
}

export async function navigateToNewRecipePage(page: Page) {
    await clickButtonByTitle(page, "Neues Rezept");
}

export async function navigateToEditPage(page: Page) {
    await clickButtonByTitle(page, "Rezept bearbeiten");
}

export async function navigateToAdminPage(page: Page) {
    await clickButtonByTitle(page, "Admin");
}

export async function navigateToHomePage(page: Page) {
    await clickButtonByTitle(page, "Home");
}

export async function fillInputByLabel(page: Page, label: string, value: string) {
    const input = page.getByLabel(label, { exact: true });
    await expect(input).toBeVisible();
    await input.fill(value)
}

export async function fillInputByPlaceholder(page: Page, placeholder: string, value: string) {
    const input = page.getByPlaceholder(placeholder, { exact: true });
    await expect(input).toBeVisible();
    await input.fill(value)
}

export async function submitForm(page: Page, name: string) {
    const button = page.getByRole("button", { name, exact: true });
    await expect(button).toBeVisible();
    await button.click();
}

export async function assertTextVisible(page: Page, text: string) {
    await expect(
        page.locator("#content").getByText(text),
    ).toBeVisible();
}

