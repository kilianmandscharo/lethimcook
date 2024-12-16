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
    ).then(res => res.json()).then(recipes => {
        if (recipes.length === 0) {
            return;
        }
        const insertLinkIntoTextfield = (title, id) => {
            const link = `[${title}](${window.location.origin}/recipe/${id})`;
            const newCursorPos = cursorPos + link.length;
            target.value = target.value.slice(0, i) + link + target.value.slice(cursorPos);
            target.selectionStart = target.selectionEnd = newCursorPos - (cursorPos - i);
        }
        if (recipes.length === 1) {
            const recipe = recipes[0];
            insertLinkIntoTextfield(recipe.title, recipe.id);
        } else {
            openSelectDialog(recipes, insertLinkIntoTextfield);
        }
    }).catch(console.error);
}

function openSelectDialog(options, cb) {
    const dialogContainer = document.createElement("div");
    dialogContainer.style.top = dialogContainer.getBoundingClientRect().top;
    dialogContainer.id = "select-dialog-container";

    const closeDialog = () => {
        document.getElementById("select-dialog-container")?.remove();
    }

    const listenForCloseDialog = () => {
        closeDialog();
        document.removeEventListener("keydown", listenForCloseDialog);
    }

    document.addEventListener("keydown", listenForCloseDialog);

    const dialog = document.createElement("div");
    dialog.id = "select-dialog";

    dialogContainer.appendChild(dialog);

    const dialogHeader = document.createElement("div");
    dialogHeader.className = "header";

    const dialogHeaderTitleContainer = document.createElement("div");
    dialogHeaderTitleContainer.className = "title-container";

    const dialogHeaderTitle = document.createElement("div");
    dialogHeaderTitle.innerText = "Rezept auswählen";

    const titleIcon = document.createElement("i");
    titleIcon.className = "fa-solid fa-caret-right";

    dialogHeaderTitleContainer.appendChild(titleIcon)
    dialogHeaderTitleContainer.appendChild(dialogHeaderTitle)
    dialogHeader.appendChild(dialogHeaderTitleContainer);

    const icon = document.createElement("i");
    icon.className = "fa-regular fa-circle-xmark fa-xl";
    icon.onclick = closeDialog;
    dialogHeader.appendChild(icon)

    const optionsContainer = document.createElement("div");
    optionsContainer.id = "options-container";

    options.forEach(option => {
        const optionElement = document.createElement("div");
        optionElement.innerText = option.title;
        optionElement.onclick = () => {
            cb(option.title, option.id);
            closeDialog();
        }
        optionsContainer.appendChild(optionElement);
    })

    const dialogFooter = document.createElement("div");

    dialog.appendChild(dialogHeader);
    dialog.appendChild(optionsContainer);
    dialog.appendChild(dialogFooter);

    document.body.appendChild(dialogContainer);
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
