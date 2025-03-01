export function attachMutationObserverToNotificationContainer() {
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

export function notify(message: string, type: "success" | "error") {
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

