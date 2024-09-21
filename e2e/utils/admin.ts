import { Page, expect } from "@playwright/test";
import { fillInputByPlaceholder, submitForm, clickButtonByTitle, assertTextVisible } from "./utils";

export async function login(page: Page, password: string) {
    await fillInputByPlaceholder(page, "Passwort", password);
    await submitForm(page, "Anmelden");
    await expect(page.getByRole("button", { name: "Abmelden" })).toBeVisible();
    await expect(
        page.getByPlaceholder("Passwort", { exact: true }),
    ).toBeDisabled();
    await expect(page.getByRole("button", { name: "Anmelden" })).toBeDisabled();
}

export async function logout(page: Page) {
    await clickButtonByTitle(page, "Abmelden");
    await expect(
        page.getByRole("button", { name: "Abmelden" }),
    ).not.toBeVisible();
    await expect(
        page.getByPlaceholder("Passwort", { exact: true }),
    ).toBeEnabled();
    await expect(page.getByRole("button", { name: "Anmelden" })).toBeEnabled();
}

export async function changePassword(
    page: Page,
    oldPassword: string,
    newPassword: string,
) {
    await fillInputByPlaceholder(page, "Altes Passwort", oldPassword);
    await fillInputByPlaceholder(page, "Neues Passwort", newPassword);
    await submitForm(page, "Bestätigen");
    await expect(page.getByPlaceholder("Altes Passwort")).toBeEmpty();
    await expect(page.getByPlaceholder("Neues Passwort")).toBeEmpty();
}

export async function testInvalidPasswordChange(page: Page) {
    await changePassword(page, "invalid", "nimda");
    await assertTextVisible(page, "Falsches Passwort")
}

export async function testPasswordChangeTooShort(page: Page) {
    await changePassword(page, "admin", "aaa");
    await assertTextVisible(page, "Minimale Passwortlänge: 5")
}

export async function testInvalidPassword(page: Page, invalidPassword: string) {
    await fillInputByPlaceholder(page, "Passwort", invalidPassword);
    await submitForm(page, "Anmelden");
    await page.waitForLoadState("networkidle");
    await expect(
        page.getByRole("button", { name: "Abmelden" }),
    ).not.toBeVisible();
    await assertTextVisible(page, "Falsches Passwort");
}

