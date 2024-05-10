<script lang="ts" setup>
import { computed } from 'vue'
import { type RouteLocationRaw, RouterLink } from 'vue-router'
import BaseIcon from '@/components/UI/BaseIcon.vue'

interface Props {
  variant?: 'default' | 'second' | 'third'
  icon?: string
  to?: RouteLocationRaw
  disabled?: boolean
  appendIcon?: string
  prependIcon?: string
  type?: 'submit' | 'reset' | 'button'
  block?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  icon: '',
  to: '',
  appendIcon: '',
  prependIcon: '',
  type: 'button'
})

const isLink = computed(() => {
  return !!props.to
})

const componentButton = computed(() => {
  return isLink.value ? RouterLink : 'button'
})
</script>

<template>
  <component
    :is="componentButton"
    class="ui-button"
    :class="[{ 'ui-button_block': block }, `ui-button_${variant}`]"
    :to="to || undefined"
    :disabled="disabled"
  >
    <BaseIcon v-if="prependIcon" class="ui-button__prepend-icon" :icon="prependIcon" />
    <span v-else>
      <slot />
    </span>
    <BaseIcon v-if="appendIcon" class="ui-button__append-icon" :icon="appendIcon" size="24" />
  </component>
</template>

<style lang="scss" scoped>
.ui-button {
  padding: 12px 16px;
  border-radius: 8px;
  transition:
    var(--transition) background-color,
    var(--transition) color,
    var(--transition) border;

  @include text-14s;

  &_default {
    background-color: var(--vt-c-light-green);
    color: var(--vt-c-white);
    box-shadow: var(--shadow-green);

    &:hover {
      background-color: var(--vt-c-green);
      box-shadow: var(--shadow-black);
    }

    &:disabled {
      background-color: var(--vt-c-white-mute);
      color: var(--vt-c-black-mute);
    }
  }

  &_icon {
    padding: 13px;
  }

  &__append-icon {
    margin-left: 8px;
  }

  &_block {
    width: 100%;
  }
}
</style>
