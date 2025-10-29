<template>
  <div 
    v-if="isVisible"
    :class="alertClasses"
    role="alert"
  >
    <div class="flex items-start">
      <!-- Icon -->
      <div class="flex-shrink-0">
        <component :is="iconComponent" class="w-5 h-5" />
      </div>
      
      <!-- Content -->
      <div class="ml-3 flex-1">
        <h3 v-if="title" :class="titleClasses">
          {{ title }}
        </h3>
        <p :class="messageClasses">
          {{ message }}
        </p>
      </div>
      
      <!-- Close button -->
      <div v-if="dismissible" class="ml-auto pl-3">
        <button
          @click="dismiss"
          :class="closeButtonClasses"
          type="button"
        >
          <span class="sr-only">Fermer</span>
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, h } from 'vue'
import type { AlertProps, AlertType } from '@/types/ui'

// Props
const props = withDefaults(defineProps<AlertProps>(), {
  dismissible: false
})

// Emits
const emit = defineEmits<{
  dismiss: []
}>()

// State
const isVisible = ref<boolean>(true)

// Icons components
const CheckCircleIcon = () => h('svg', {
  fill: 'currentColor',
  viewBox: '0 0 20 20'
}, h('path', {
  'fill-rule': 'evenodd',
  d: 'M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z',
  'clip-rule': 'evenodd'
}))

const ExclamationCircleIcon = () => h('svg', {
  fill: 'currentColor',
  viewBox: '0 0 20 20'
}, h('path', {
  'fill-rule': 'evenodd',
  d: 'M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z',
  'clip-rule': 'evenodd'
}))

const ExclamationTriangleIcon = () => h('svg', {
  fill: 'currentColor',
  viewBox: '0 0 20 20'
}, h('path', {
  'fill-rule': 'evenodd',
  d: 'M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z',
  'clip-rule': 'evenodd'
}))

const InformationCircleIcon = () => h('svg', {
  fill: 'currentColor',
  viewBox: '0 0 20 20'
}, h('path', {
  'fill-rule': 'evenodd',
  d: 'M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z',
  'clip-rule': 'evenodd'
}))

// Computed
const alertClasses = computed((): string => {
  const baseClasses = 'rounded-lg p-4 border transition-all duration-200'
  
  const typeClasses: Record<AlertType, string> = {
    success: 'bg-green-50 border-green-200 text-green-800',
    error: 'bg-red-50 border-red-200 text-red-800',
    warning: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    info: 'bg-blue-50 border-blue-200 text-blue-800'
  }
  
  return `${baseClasses} ${typeClasses[props.type]}`
})

const titleClasses = computed((): string => {
  const typeClasses: Record<AlertType, string> = {
    success: 'text-green-800',
    error: 'text-red-800',
    warning: 'text-yellow-800',
    info: 'text-blue-800'
  }
  
  return `font-medium text-sm ${typeClasses[props.type]}`
})

const messageClasses = computed((): string => {
  const baseClasses = 'text-sm'
  const marginClass = props.title ? 'mt-1' : ''
  
  const typeClasses: Record<AlertType, string> = {
    success: 'text-green-700',
    error: 'text-red-700',
    warning: 'text-yellow-700',
    info: 'text-blue-700'
  }
  
  return `${baseClasses} ${marginClass} ${typeClasses[props.type]}`
})

const closeButtonClasses = computed((): string => {
  const baseClasses = 'inline-flex rounded-md p-1.5 focus:outline-none focus:ring-2 focus:ring-offset-2'
  
  const typeClasses: Record<AlertType, string> = {
    success: 'text-green-500 hover:bg-green-100 focus:ring-green-600',
    error: 'text-red-500 hover:bg-red-100 focus:ring-red-600',
    warning: 'text-yellow-500 hover:bg-yellow-100 focus:ring-yellow-600',
    info: 'text-blue-500 hover:bg-blue-100 focus:ring-blue-600'
  }
  
  return `${baseClasses} ${typeClasses[props.type]}`
})

const iconComponent = computed(() => {
  const icons = {
    success: CheckCircleIcon,
    error: ExclamationCircleIcon,
    warning: ExclamationTriangleIcon,
    info: InformationCircleIcon
  }
  
  return icons[props.type]
})

// Methods
const dismiss = (): void => {
  isVisible.value = false
  emit('dismiss')
}
</script>
