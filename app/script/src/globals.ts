import { notify } from "./notification";
import {
    fetchLink,
    getLastListNumberBeforeCursor,
    insertStringInInputAtCursor,
} from "./utils";

export function attachTextAreaEventListeners() {
    const handleIngredientsKeyup = (e: KeyboardEvent) => {
        if (e.key === "Enter" && e.target) {
            insertStringInInputAtCursor(e.target as HTMLTextAreaElement, "- ");
        }
    };
    const handleIngredientsInput = (e: Event) => {
        if (e.target) {
            fetchLink(e.target as HTMLTextAreaElement);
        }
    };
    const handleInstructionsKeyup = (e: KeyboardEvent) => {
        if (e.key === "Enter" && e.target) {
            const lastListNumber = getLastListNumberBeforeCursor(
                e.target as HTMLTextAreaElement,
            );
            if (lastListNumber) {
                insertStringInInputAtCursor(
                    e.target as HTMLTextAreaElement,
                    `${parseInt(lastListNumber) + 1}. `,
                );
            }
        }
    };
    const handleInstructionsInput = (e: Event) => {
        if (e.target) {
            fetchLink(e.target as HTMLTextAreaElement);
        }
    };

    const ingredients = document.getElementById("ingredients");
    ingredients?.addEventListener("keyup", handleIngredientsKeyup);
    ingredients?.addEventListener("input", handleIngredientsInput);

    const instructions = document.getElementById("instructions");
    instructions?.addEventListener("keyup", handleInstructionsKeyup);
    instructions?.addEventListener("input", handleInstructionsInput);
}

export function handleCopyUrlToClipboardButton() {
    const url = window.location.href;
    navigator.clipboard
        .writeText(url)
        .then(() => {
            notify("Link kopiert", "success");
        })
        .catch(() => {
            notify("Kopieren fehlgeschlagen", "error");
        });
}
