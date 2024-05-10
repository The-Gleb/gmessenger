<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import BaseInput from './BaseInput.vue'
import BaseIcon from './BaseIcon.vue'
import { type InputProps, type InputEmits } from './BaseInput'

const props = defineProps<InputProps>()
const emit = defineEmits<InputEmits>()

const updatedValue = defineModel<string>({ required: true })

const isInputValueEmpty = computed(() => {
  return !updatedValue.value.length
})

const passwordShowState = ref(false)
const passwordShowIcon = computed(() => {
  return passwordShowState.value ? 'password-visible' : 'password-hidden'
})

const inputType = computed(() => {
  return passwordShowState.value ? 'text' : 'password'
})

watch(isInputValueEmpty, () => {
  passwordShowState.value = false
})
</script>

<template>
  <BaseInput
    v-bind="props"
    v-model="updatedValue"
    :type="inputType"
    @blur="emit('blur', $event)"
    @focus="emit('focus', $event)"
  >
    <template #append>
      <button
        v-if="!isInputValueEmpty"
        type="button"
        @click="passwordShowState = !passwordShowState"
      >
        <BaseIcon :icon="passwordShowIcon" size="20" color="#fff" />
      </button>
    </template>
  </BaseInput>
</template>
