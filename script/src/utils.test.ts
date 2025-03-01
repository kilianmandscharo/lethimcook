import { extractContentBetweenBrackets } from "./utils";

describe("utils", () => {
    describe("extractContentBetweenBrackets", () => {
        const testCases = [
            {
                target: {
                    value: "[!Naan]",
                    selectionStart: 7,
                },
                expected: {
                    query: "Naan",
                    substitutionStart: 0,
                    substitutionEnd: 7,
                },
            },
            {
                target: {
                    value: "some text before [!Naan]",
                    selectionStart: 24,
                },
                expected: {
                    query: "Naan",
                    substitutionStart: 17,
                    substitutionEnd: 24,
                },
            },
            {
                target: {
                    value: "some text before [!Naan] some text after",
                    selectionStart: 24,
                },
                expected: {
                    query: "Naan",
                    substitutionStart: 17,
                    substitutionEnd: 24,
                },
            },
            {
                target: {
                    value: "[Naan]",
                    selectionStart: 6,
                },
                expected: null,
            },
            {
                target: {
                    value: "!Naan]",
                    selectionStart: 6,
                },
                expected: null,
            },
            {
                target: {
                    value: "",
                    selectionStart: 1,
                },
                expected: null,
            },
        ];

        for (const test of testCases) {
            it(`should extract ${test.target.value} correctly`, () => {
                const result = extractContentBetweenBrackets(test.target);
                expect(result).toStrictEqual(test.expected);
            });
        }
    });
});
