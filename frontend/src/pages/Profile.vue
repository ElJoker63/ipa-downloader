<template>
  <div class="max-w-4xl mx-auto space-y-8 animate-slide-up font-sans">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-white flex items-center space-x-3">
          <span>{{ t.profile.title }}</span>
          <span class="px-2.5 py-0.5 text-[10px] font-bold rounded-full bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30 uppercase">
            {{ t.profile.activeSession }}
          </span>
        </h1>
        <p class="text-sm text-[#8E8E93] mt-1">{{ t.profile.subtitle }}</p>
      </div>

      <button
        @click="logout"
        class="px-6 py-2.5 rounded-2xl bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white border border-red-500/30 transition-all font-bold text-sm"
      >
        {{ t.profile.signOut }}
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-12 gap-8">
      <!-- Account Info Card -->
      <div class="md:col-span-4 flex flex-col items-center space-y-6">
        <div class="w-32 h-32 rounded-[32px] bg-gradient-to-tr from-[#0071E3] via-[#0A84FF] to-[#64D2FF] flex items-center justify-center text-4xl font-bold text-white shadow-2xl shadow-[#0A84FF]/40 border-4 border-white/10 shrink-0">
          {{ userInitials }}
        </div>
        <div class="text-center">
          <h2 class="text-xl font-bold text-white">{{ authStore.account.name || 'Apple User' }}</h2>
          <p class="text-sm text-[#8E8E93] truncate max-w-[240px]">{{ authStore.account.email }}</p>
        </div>
      </div>

      <!-- Detail Grid -->
      <div class="md:col-span-8 space-y-6">
        <div class="glass-card p-6 rounded-[24px] space-y-6">
          <h3 class="text-xs font-bold uppercase tracking-widest text-[#7D8592]">{{ t.profile.sessionMetadata }}</h3>

          <div class="grid grid-cols-2 gap-6">
            <div class="space-y-1.5">
              <span class="text-[10px] uppercase font-bold text-[#7D8592]">{{ t.profile.storefrontId }}</span>
              <div class="text-sm text-white font-mono bg-white/[0.04] px-3 py-2 rounded-xl border border-white/[0.08]">
                {{ authStore.account.storeFront || '--' }}
              </div>
            </div>
            <div class="space-y-1.5">
              <span class="text-[10px] uppercase font-bold text-[#7D8592]">{{ t.profile.region }}</span>
              <div class="text-sm text-white bg-white/[0.04] px-3 py-2 rounded-xl border border-white/[0.08] flex items-center space-x-2">
                <span>{{ authStore.account.storeFrontCountry || 'Unknown' }}</span>
              </div>
            </div>
            <div class="space-y-1.5">
              <span class="text-[10px] uppercase font-bold text-[#7D8592]">{{ t.profile.dsid }}</span>
              <div class="text-sm text-white font-mono bg-white/[0.04] px-3 py-2 rounded-xl border border-white/[0.08]">
                {{ authStore.account.directoryServicesId || '--' }}
              </div>
            </div>
            <div class="space-y-1.5">
              <span class="text-[10px] uppercase font-bold text-[#7D8592]">{{ t.profile.accountPod }}</span>
              <div class="text-sm text-white font-mono bg-white/[0.04] px-3 py-2 rounded-xl border border-white/[0.08]">
                {{ authStore.account.pod || 'N/A' }}
              </div>
            </div>
          </div>
        </div>

        <div class="glass-card p-6 rounded-[24px] space-y-4">
           <h3 class="text-xs font-bold uppercase tracking-widest text-[#7D8592]">{{ t.profile.securityPrivacy }}</h3>
           <div class="flex items-center justify-between p-4 rounded-2xl bg-[#30D158]/5 border border-[#30D158]/20">
              <div class="flex items-center space-x-3">
                 <div class="w-10 h-10 rounded-xl bg-[#30D158]/10 flex items-center justify-center text-[#30D158]">
                   <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" /></svg>
                 </div>
                 <div>
                   <div class="text-sm font-bold text-white">{{ t.profile.keychainTitle }}</div>
                   <p class="text-[11px] text-[#B8C0CC]">{{ t.profile.keychainDesc }}</p>
                 </div>
              </div>
              <span class="text-[10px] font-bold text-[#30D158] uppercase">{{ t.profile.activeStatus }}</span>
           </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useModalStore } from '../stores/modal'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'

const router = useRouter()
const authStore = useAuthStore()
const modalStore = useModalStore()
const { t } = useI18n()
const { showToast } = useNotifications()

const userInitials = computed(() => {
  const name = authStore.account?.name || authStore.account?.email || ''
  if (!name.trim()) return ''
  const parts = name.trim().split(/\s+/)
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase()
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
})

async function logout() {
  modalStore.confirm(
    t.value.common.confirm,
    t.value.profile.signOut + '?',
    async () => {
      await authStore.logout()
      showToast(t.value.auth.signedOut, t.value.auth.sessionRevoked, 'info')
      router.push('/')
    },
    t.value.profile.signOut
  )
}
</script>
