const INPUTS = ["title", "description", "duration", "author", "source", "tags", "ingredients", "instructions"];
const RECIPE_NEW_KEY = "recipe_new";

function main() {
    attachMutationObserverToNotificationContainer();
}

function attachMutationObserverToNotificationContainer() {
    const notificationContainer = document.getElementById("notification-container");
    if (!notificationContainer) {
        return;
    }
    const observer = new MutationObserver((mutationsList, _) => {
        for (let mutation of mutationsList) {
            if (mutation.type === 'childList') {
                mutation.addedNodes.forEach(node => {
                    setTimeout(() => {
                        notificationContainer.removeChild(node);
                    }, 3000);
                });
            }
        }
    });
    observer.observe(notificationContainer, { childList: true });
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
            })
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
};

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
    newNotification.className =
        type === "success"
            ? "notification notification-success"
            : "notification notification-danger";

    const newNotificationIcon = document.createElement("div");
    newNotificationIcon.className =
        type === "success"
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
