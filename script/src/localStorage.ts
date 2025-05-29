export class LocalStorageUtil {
    static INPUTS = [
        "title",
        "description",
        "cookingDuration",
        "totalDuration",
        "author",
        "source",
        "tags",
        "ingredients",
        "instructions",
    ];

    static RECIPE_NEW_KEY = "recipe_new";

    static saveForm() {
        const recipe = LocalStorageUtil.INPUTS.reduce(
            (acc, inputName) => {
                const input = document.getElementById(inputName) as
                    | HTMLInputElement
                    | HTMLTextAreaElement;
                acc[inputName] = input?.value ?? "";
                return acc;
            },
            {} as { [key: string]: string },
        );
        localStorage.setItem(
            LocalStorageUtil.RECIPE_NEW_KEY,
            JSON.stringify(recipe),
        );
    }

    static loadForm() {
        const recipeString = localStorage.getItem(
            LocalStorageUtil.RECIPE_NEW_KEY,
        );
        if (!recipeString) {
            return;
        }
        const recipe = JSON.parse(recipeString);
        if (!recipe) {
            return;
        }
        LocalStorageUtil.INPUTS.forEach((inputName) => {
            const input = document.getElementById(inputName) as
                | HTMLInputElement
                | HTMLTextAreaElement;
            if (input) {
                input.value = recipe[inputName] ?? "";
            }
        });
    }

    static deleteForm() {
        localStorage.setItem(LocalStorageUtil.RECIPE_NEW_KEY, "");
    }
}
