<template>
  <div class="relative inline-block" ref="dropdownRef">
    <!-- Trigger -->
    <button
      type="button"
      @click="isOpen = !isOpen"
      class="flex items-center justify-between w-full bg-white/[0.06] border border-white/[0.1] hover:border-white/[0.2] rounded-xl px-3.5 py-2 text-xs text-white transition-all duration-200 active:scale-[0.98]"
      :class="[isOpen ? 'border-[#0A84FF]/50 ring-2 ring-[#0A84FF]/20' : '']"
    >
      <div class="flex items-center space-x-2 truncate mr-2">
        <slot name="icon" :selected="selectedOption"></slot>
        <span class="truncate">{{ selectedLabel }}</span>
      </div>
      <svg
        class="w-3.5 h-3.5 text-[#7D8592] transition-transform duration-200 shrink-0"
        :class="{ 'rotate-180': isOpen }"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- Dropdown Menu -->
    <Transition name="dropdown-fade">
      <div
        v-if="isOpen"
        class="absolute z-[100] mt-2 w-full min-w-[160px] bg-[#1C1C1E]/95 backdrop-blur-2xl border border-white/10 rounded-2xl shadow-2xl overflow-hidden py-1.5"
      >
        <div class="max-h-60 overflow-y-auto custom-scrollbar">
          <button
            v-for="option in options"
            :key="option.id"
            type="button"
            @click="selectOption(option)"
            class="w-full px-3 py-2 text-left text-xs transition-colors flex items-center space-x-2"
            :class="[
              modelValue === option.id
                ? 'bg-[#0A84FF] text-white font-bold'
                : 'text-[#B8C0CC] hover:bg-white/[0.08] hover:text-white'
            ]"
          >
            <span v-if="option.flag" class="text-sm leading-none">{{ option.flag }}</span>
            <span class="truncate">{{ option.name }}</span>
            <svg v-if="modelValue === option.id" class="w-3.5 h-3.5 ml-auto text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Option {
  id: string | number
  name: string
  flag?: string
}

const props = defineProps<{
  modelValue: string | number
  options: Option[]
  placeholder?: string
}>()

const emit = defineEmits(['update:modelValue', 'change'])

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

const selectedOption = computed(() => {
  return props.options.find(o => o.id === props.modelValue)
})

const selectedLabel = computed(() => {
  return selectedOption.value?.name || props.placeholder || 'Select...'
})

function selectOption(option: Option) {
  emit('update:modelValue', option.id)
  emit('change', option.id)
  isOpen.value = false
}

// Close when clicking outside
const handleClickOutside = (event: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside)
})
</script>

<style scoped>
.dropdown-fade-enter-active, .dropdown-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.dropdown-fade-enter-from, .dropdown-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.98);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
