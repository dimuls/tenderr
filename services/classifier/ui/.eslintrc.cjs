module.exports = {
    env: {
        node: true,
    },
    extends: ["eslint:recommended", "plugin:vue/vue3-recommended", "prettier"],
    rules: {
        "vue/attributes-order": "off",
        "vue/multi-word-component-names": "off",
        "vue/order-in-components": "off",
    },
};