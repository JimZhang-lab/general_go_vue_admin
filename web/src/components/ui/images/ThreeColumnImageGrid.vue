/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/components/ui/images/ThreeColumnImageGrid.vue
 * @Description: 
 * 
 */
<template>
  <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-3">
    <div v-for="(image, index) in images" :key="index">
      <img
        :src="image.src"
        :alt="image.alt"
        class="w-full border border-gray-200 rounded-xl dark:border-gray-800"
      />
    </div>
  </div>
</template>

<script setup>
const images = [
  { src: '/images/grid-image/image-04.png', alt: 'Grid image 1' },
  { src: '/images/grid-image/image-05.png', alt: 'Grid image 2' },
  { src: '/images/grid-image/image-06.png', alt: 'Grid image 3' },
]
</script>
