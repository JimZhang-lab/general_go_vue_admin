/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/components/ui/images/TwoColumnImageGrid.vue
 * @Description: 
 * 
 */
<template>
  <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
    <div>
      <img
        src="/images/grid-image/image-02.png"
        alt="grid"
        class="w-full border border-gray-200 rounded-xl dark:border-gray-800"
      />
    </div>

    <div>
      <img
        src="/images/grid-image/image-03.png"
        alt="grid"
        class="w-full border border-gray-200 rounded-xl dark:border-gray-800"
      />
    </div>
  </div>
</template>
