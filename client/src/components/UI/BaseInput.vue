<script lang="ts" setup>
import { computed, onMounted, ref, useSlots } from 'vue'
import BaseIcon from '@/components/UI/BaseIcon.vue'
import { type InputProps, type InputEmits } from './BaseInput'

const props = withDefaults(defineProps<InputProps>(), {
  prependIcon: undefined,
  prependIconColor: 'text-secondary',
  appendIcon: undefined,
  appendIconColor: 'text-secondary',
  label: undefined,
  hint: undefined,
  density: 'comfortable',
  type: 'text',
  placeholder: '',
  autocomplete: undefined,
  errorMessages: undefined,
  required: false
})
const emit = defineEmits<InputEmits>()

const updatedValue = defineModel<string>({ required: true })

const inputRef = ref<HTMLInputElement | null>(null)
const inputLabel = ref<HTMLLabelElement | null>(null)

onMounted(() => {
  if (props.autofocus) {
    inputLabel.value?.focus()
  }
})

const slots = useSlots()
const hasPrependSlot = computed(() => slots.prepend)
const hasAppendSlot = computed(() => slots.append)

const getIconTag = (state: boolean) => (state ? 'button' : 'div')
const prependIconTag = computed(() => getIconTag(props.prependIconClickable))
const appendIconTag = computed(() => getIconTag(props.appendIconClickable))
</script>

<template>
  <div
    class="ui-input"
    :class="[
      {
        'ui-input_slot-prepend': hasPrependSlot,
        'ui-input_slot-append': hasAppendSlot,
        'ui-input_disabled': disabled,
        'ui-input_error': errorMessages?.length
      }
    ]"
  >
    <div v-if="label" class="ui-input__title">
      {{ label }}
      <span v-if="required">*</span>
    </div>
    <label class="ui-input__inner">
      <slot name="prepend">
        <component
          :is="prependIconTag"
          v-if="prependIcon"
          class="ui-input__icon ui-input__icon_prepend"
          :style="{ '--color': `var(--${prependIconColor})` }"
          @click.stop="emit('click:prepend-icon', $event)"
        >
          <BaseIcon :icon="prependIcon" size="16" />
        </component>
      </slot>
      <div class="ui-input__field">
        <input
          ref="inputRef"
          v-model="updatedValue"
          :type="type"
          :placeholder="placeholder"
          :autocomplete="autocomplete"
          @focus="emit('focus', $event)"
          @blur="emit('blur', $event)"
        />
      </div>
      <div class="ui-input__icons_append">
        <slot name="append">
          <component
            :is="appendIconTag"
            v-if="appendIcon"
            class="ui-input__icon ui-input__icon_append"
            :style="{ '--color': `var(--${appendIconColor})` }"
            @click.stop="emit('click:append-icon', $event)"
          >
            <BaseIcon :icon="appendIcon" size="16" />
          </component>
        </slot>
      </div>
    </label>
  </div>
</template>

<style lang="scss" scoped>
.ui-input {
  @include text-small-13s;

  &__title {
    color: var(--vt-c-white-mute);
    margin-bottom: 4px;
    transition: var(--transition) color;

    @include text-14s;

    span {
      color: var(--vt-c-red);
    }
  }

  &__inner {
    width: 100%;
    height: 40px;
    padding: 0 8px;
    border: 1px solid var(--vt-c-white-mute);
    border-radius: 6px;
    background-color: var(--vt-c-black-mute);
    display: flex;
    align-items: center;
    cursor: text;
    transition: var(--transition) border-color;
  }

  &__field {
    width: 100%;

    > input {
      background-color: transparent;
      width: 100%;
      outline: none;
      color: var(--vt-c-white);
      line-height: 1;

      &::placeholder {
        color: var(--vt-c-white-mute);
      }
    }
  }

  &__icon {
    flex: 0 0 16px;
    color: var(--color);
  }

  &__icons {
    &_append {
      padding-left: 4px;
      display: flex;
      align-items: center;

      > * + * {
        margin-left: 12px;
      }
    }
  }

  &_disabled {
    .ui-input__inner {
      background-color: var(--vt-c-white-mute);
    }
  }

  &_error {
    .ui-input__title {
      color: var(--vt-c-red);
    }

    .ui-input__inner {
      border-color: var(--vt-c-red);
    }

    .ui-input__details {
      color: var(--vt-c-red);
    }
  }
}
</style>
