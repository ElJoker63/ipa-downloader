# Implementation Plan: Proper Light Theme Support

The current project has partial light theme support, but it's inconsistent because many UI components use hardcoded hex colors that only look good in dark mode. Additionally, the Tailwind configuration is hardcoded to dark mode hex values, making it difficult to use Tailwind's utility classes for themed components.

## Proposed Changes

### 1. Style Foundations

#### [MODIFY] [tailwind.config.js](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/tailwind.config.js)
- Update the `colors` section to reference CSS variables instead of hex values. This allows `bg-bg-primary`, `text-text-secondary`, etc., to automatically switch values based on the active theme.

#### [MODIFY] [index.css](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/index.css)
- Refine the CSS variables in `:root` (dark) and `html.light`.
- Add status colors to `html.light` to ensure they have proper contrast.
- Update global component styles (`.glass-panel`, `.btn-primary`, etc.) to use the variables consistently.

### 2. Component Refactoring

I will systematically replace hardcoded hex colors with the corresponding Tailwind semantic classes.

#### [MODIFY] [MainLayout.vue](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/layouts/MainLayout.vue)
- Replace `bg-[#0F1115]` with `bg-bg-primary`.
- Replace `text-[#FFFFFF]` with `text-text-primary`.
- Update sidebar navigation colors.

#### [MODIFY] [Apps.vue](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/pages/Apps.vue)
- Update device selection cards and queue items.
- Replace hardcoded border and background colors.

#### [MODIFY] [AppDetailsModal.vue](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/components/AppDetailsModal.vue)
- Update modal backgrounds, text colors, and metadata labels.

#### [MODIFY] [Downloads.vue](file:///D:/_PROJETCS/Github/ipa-downloader/frontend/src/pages/Downloads.vue)
- Update progress bars and task items.

#### [MODIFY] [Other Components]
- Similar updates to `Settings.vue`, `Favorites.vue`, `ToastContainer.vue`, and other shared components.

## Verification Plan

### Manual Verification
- Toggle between Light, Dark, and System themes in the Settings page.
- Inspect all main views (Apps, Downloads, Favorites, Settings) to ensure no hardcoded dark elements remain.
- Verify that modals and overlays adapt correctly to the light theme.
- Check that text remains legible across both themes.
