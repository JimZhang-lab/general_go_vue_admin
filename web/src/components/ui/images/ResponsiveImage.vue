/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/components/ui/images/ResponsiveImage.vue
 * @Description: 
 * 
 */
<template>
  <div class="relative">
    <div id="pane" class="overflow-hidden">
      <img
        src="/images/grid-image/image-01.png"
        alt="Cover"
        class="w-full border border-gray-200 rounded-xl dark:border-gray-800"
      />
    </div>
    <div id="ghostpane" class="absolute top-0 left-0 duration-300 ease-in-out"></div>
  </div>
</template>
