# Reversion Plan: Restore Original Dark-Only Design

The user has requested to remove the light theme support and restore the original design which was focused on Dark Mode.

## Proposed Changes

### 1. Restore Foundations

#### [MODIFY] [tailwind.config.js](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/tailwind.config.js)
- Restore hardcoded hex values in the `colors` section.

#### [MODIFY] [index.css](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/index.css)
- Remove `html.light` overrides.
- Restore original hex values for CSS variables.

### 2. Restore Component Styles

#### [MODIFY] [MainLayout.vue](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/layouts/MainLayout.vue)
#### [MODIFY] [Apps.vue](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/pages/Apps.vue)
#### [MODIFY] [AppDetailsModal.vue](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/components/AppDetailsModal.vue)
- Revert semantic Tailwind classes back to original hardcoded hex classes (e.g., `bg-[#0F1115]`, `text-[#FFFFFF]`).

## Verification Plan
- Verify that the app looks exactly as it did before the theme refactor.
- Ensure only Dark Mode is active.
