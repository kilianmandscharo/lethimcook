const INPUTS = [
    "title",
    "description",
    "duration",
    "author",
    "source",
    "tags",
    "ingredients",
    "instructions",
];
const RECIPE_NEW_KEY = "recipe_new";

function main() {
    attachMutationObserverToNotificationContainer();
}

function attachMutationObserverToNotificationContainer() {
    const notificationContainer = document.getElementById(
        "notification-container",
    );
    if (!notificationContainer) {
        return;
    }
    const observer = new MutationObserver((mutationsList, _) => {
        for (const mutation of mutationsList) {
            if (mutation.type === "childList") {
                mutation.addedNodes.forEach((node) => {
                    setTimeout(() => {
                        notificationContainer.removeChild(node);
                    }, 3000);
                });
            }
        }
    });
    observer.observe(notificationContainer, { childList: true });
}

let state = {};

function injectLinkIntoTextarea(title, id) {
    const link = `[${title}](${window.location.origin}/recipe/${id})`;
    const newCursorPos = state.cursorPos + link.length;
    state.target.value = state.target.value.slice(0, state.i) + link + state.target.value.slice(state.cursorPos);
    state.target.selectionStart = state.target.selectionEnd = newCursorPos - (state.cursorPos - state.i);
    state = {};
}

function fetchLink(target) {
    const cursorPos = target.selectionStart;
    if (target.value.length === 0) {
        return;
    }
    const lastChar = target.value[cursorPos - 1];
    if (lastChar !== "]") {
        return;
    }
    let bracketContentReversed = "";
    let i = cursorPos - 2;
    while (i-- > 0) {
        if (target.value[i] === "[") {
            break;
        }
        bracketContentReversed += target.value[i];
    }
    const params = new URLSearchParams({
        query: bracketContentReversed.split("").reverse().join(""),
    });
    fetch(
        `${window.location.origin}/recipe/link?` + params.toString()
    ).then(res => res.text()).then(data => {
        state.target = target;
        state.cursorPos = cursorPos;
        state.i = i;
        try {
            const recipe = JSON.parse(data);
            injectLinkIntoTextarea(recipe.title, recipe.id);
        } catch {
            const container = document.createElement("div");
            container.id = "select-dialog-container";
            container.style.top = container.getBoundingClientRect().top;
            container.innerHTML = data;
            document.body.appendChild(container);
            setTimeout(() => {
                const closeDialog = () => {
                    document.getElementById("select-dialog-container")?.remove();
                }
                const listenForCloseDialog = () => {
                    closeDialog();
                    document.removeEventListener("keydown", listenForCloseDialog);
                }
                document.addEventListener("keydown", listenForCloseDialog);
            }, 0);
        }
    }).catch(console.error);
}

function attachTextAreaKeyupEventListeners() {
    const ingredients = document.getElementById("ingredients");
    if (ingredients) {
        ingredients.addEventListener("keyup", (e) => {
            if (e.key === "Enter") {
                insertStringInInputAtCursor(e.target, "- ");
            }
        });
        ingredients.addEventListener("input", (e) => {
            fetchLink(e.target);
        });
    }
    const instructions = document.getElementById("instructions");
    if (instructions) {
        instructions.addEventListener("keyup", (e) => {
            if (e.key === "Enter") {
                const lastListNumber = getLastListNumberBeforeCursor(e.target);
                if (lastListNumber) {
                    insertStringInInputAtCursor(
                        e.target,
                        `${parseInt(lastListNumber) + 1}. `,
                    );
                }
            }
        });
        instructions.addEventListener("input", (e) => {
            fetchLink(e.target);
        });
    }
}

function getLastListNumberBeforeCursor(el) {
    const re = /(\n(\d+)|^(\d+))\./g;
    const matches = [...el.value.slice(0, el.selectionStart).matchAll(re)];
    return matches.length > 0 ? matches[matches.length - 1][1] : null;
}

function insertStringInInputAtCursor(el, s) {
    const cursorPos = el.selectionStart;
    const newCursorPos = cursorPos + s.length;
    el.value = el.value.slice(0, cursorPos) + s +
        el.value.slice(el.selectionEnd);
    el.selectionStart = el.selectionEnd = newCursorPos;
}

function attachNewFormSubmitEventListener() {
    const submit = document.getElementById("recipe-new-submit");
    if (!submit) {
        return;
    }
    submit.addEventListener("click", () => {
        deleteFormFromLocalStorage();
    });
}

function attachFormEventListeners() {
    INPUTS.forEach((inputName) => {
        const input = document.getElementById(inputName);
        if (input) {
            input.addEventListener("blur", () => {
                saveFormToLocalStorage();
            });
        }
    });
}

function saveFormToLocalStorage() {
    const recipe = INPUTS.reduce((acc, inputName) => {
        const input = document.getElementById(inputName);
        acc[inputName] = input?.value ?? "";
        return acc;
    }, {});
    localStorage.setItem(RECIPE_NEW_KEY, JSON.stringify(recipe));
}

function loadFormFromLocalStorage() {
    const recipeString = localStorage.getItem(RECIPE_NEW_KEY);
    if (!recipeString) {
        return;
    }
    const recipe = JSON.parse(recipeString);
    if (!recipe) {
        return;
    }
    INPUTS.forEach((inputName) => {
        const input = document.getElementById(inputName);
        if (input) {
            input.value = recipe[inputName] ?? "";
        }
    });
}

function deleteFormFromLocalStorage() {
    localStorage.setItem(RECIPE_NEW_KEY, "");
}

function attachEventListenerToClipboardButton() {
    const button = document.getElementById("copy-url-to-clipboard-button");
    if (!button) {
        return;
    }
    button.addEventListener("click", () => {
        const url = window.location.href;
        navigator.clipboard.writeText(url).then(() => {
            notify("Link kopiert", "success");
        }).catch(() => {
            notify("Kopieren fehlgeschlagen", "error");
        });
    });
}

function notify(message, type) {
    const notificationContainer = document.getElementById(
        "notification-container",
    );
    if (!notificationContainer) {
        return;
    }

    const notifications = document.querySelectorAll(".notification");
    const index = notifications.length;

    const newNotification = document.createElement("div");
    newNotification.id = `notification-${index}`;
    newNotification.className = type === "success"
        ? "notification notification-success"
        : "notification notification-danger";

    const newNotificationIcon = document.createElement("div");
    newNotificationIcon.className = type === "success"
        ? "fa-solid fa-circle-info"
        : "fa-regular fa-circle-exclamation";

    const newNotificationText = document.createElement("p");
    newNotificationText.className = "notification-text";
    newNotificationText.innerText = message;

    newNotification.appendChild(newNotificationIcon);
    newNotification.appendChild(newNotificationText);

    notificationContainer.appendChild(newNotification);
}

main();
