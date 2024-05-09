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

const isOnlyIcon = computed(() => {
  return !!props.icon
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
    :class="[
      {
        'ui-button_icon': isOnlyIcon,
        'ui-button_btn': !isOnlyIcon,
        'ui-button_block': block
      },
      `ui-button_${variant}`
    ]"
    :to="to || undefined"
    :disabled="disabled"
  >
    <BaseIcon v-if="prependIcon" class="ui-button__prepend-icon" :icon="prependIcon" />
    <BaseIcon v-if="isOnlyIcon" :icon="icon" />
    <span v-else>
      <slot />
    </span>
    <BaseIcon v-if="appendIcon" class="ui-button__append-icon" :icon="appendIcon" size="24" />
  </component>
</template>

<style lang="scss" scoped>
.ui-button {
  padding: 12px 16px;
  border-radius: 32px;
  transition:
    var(--transition) background-color,
    var(--transition) color,
    var(--transition) border;

  @include text-14s;

  &_default {
    background-color: var(--icons-red);
    box-shadow: var(--shadow-red);
    color: var(--bg-default);

    &:hover {
      background-color: #e46e63;
    }

    &:active {
      background-color: var(--text-tapped);
    }

    &:disabled {
      background-color: var(--text-disabled);
      color: var(--text-tertiary);
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
