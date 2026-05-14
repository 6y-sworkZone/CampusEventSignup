<template>
  <div style="border: 1px solid #dcdfe6; border-radius: 4px; overflow: hidden;">
    <Toolbar
      style="border-bottom: 1px solid #dcdfe6;"
      :editor="editorRef"
      :defaultConfig="toolbarConfig"
      mode="default"
    />
    <Editor
      style="height: 400px; overflow-y: hidden;"
      v-model="valueHtml"
      :defaultConfig="editorConfig"
      mode="default"
      @onCreated="handleCreated"
    />
  </div>
</template>

<script setup>
import { onBeforeUnmount, ref, shallowRef, onMounted, watch } from 'vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import '@wangeditor/editor/dist/css/style.css'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = shallowRef()
const valueHtml = ref('')

onMounted(() => {
  valueHtml.value = props.modelValue
})

const toolbarConfig = {}
const editorConfig = { placeholder: '请输入活动描述...' }

const handleCreated = (editor) => {
  editorRef.value = editor
}

const handleChange = (editor) => {
  emit('update:modelValue', editor.getHtml())
}

onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

watch(valueHtml, (newVal) => {
  emit('update:modelValue', newVal)
})
</script>

<style scoped>
</style>
