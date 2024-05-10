<script setup lang="ts">
import type { ErrorObject } from '@vuelidate/core'
import { computed, nextTick, ref, unref, watch } from 'vue'

type Props = {
  errorMessages?: ErrorObject[]
}

const props = defineProps<Props>()

const inputDetailsRef = ref<HTMLDivElement | null>(null)
const inputDetailsHeight = ref(0)

const displayedDetail = computed(() => {
  const maybeError = props.errorMessages?.[0]?.$message

  if (!maybeError) {
    return
  }

  return unref(maybeError)
})

watch(displayedDetail, async () => {
  if (displayedDetail.value) {
    await nextTick()
    inputDetailsHeight.value = inputDetailsRef.value?.scrollHeight || 0
  } else {
    inputDetailsHeight.value = 0
  }
})
</script>

<template>
  <div
    class="ui-input-details"
    :class="{ 'ui-input-details_active': !!displayedDetail }"
    :style="{ '--max-height': `${inputDetailsHeight}px` }"
  >
    <div ref="inputDetailsRef">
      {{ displayedDetail }}
    </div>
  </div>
</template>

<style scoped lang="scss">
.ui-input-details {
  height: var(--max-height);
  color: var(--vt-c-red);
  font-weight: 600;
  transition:
    var(--transition) margin,
    var(--transition) height;

  &_active {
    margin-top: 8px;
  }
}
</style>
