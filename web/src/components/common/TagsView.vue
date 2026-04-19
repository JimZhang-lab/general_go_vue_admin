/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/components/common/TagsView.vue
 * @Description:
 *
 */
<template>
  <div class="tags-view-container border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 w-full shadow-sm">
    <div class="tags-view-wrapper flex overflow-x-auto no-scrollbar px-5 py-1" ref="scrollPane">
      <router-link
        v-for="tag in visitedViews"
        ref="tag"
        :key="tag.path"
        :to="{ path: tag.path, query: tag.query }"
        :class="['tags-view-item relative inline-flex items-center px-3 py-1.5 mx-1 border text-sm font-medium rounded-md transition-colors duration-200 whitespace-nowrap', isActive(tag) ? 'bg-brand-50 border-brand-200 text-brand-600 dark:bg-brand-500/10 dark:border-brand-500/30 dark:text-brand-400' : 'bg-white border-gray-200 text-gray-500 hover:bg-gray-50 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700']"
        @click.middle="!isAffix(tag) ? closeSelectedTag(tag) : ''"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        <span class="mr-1">{{ tag.title }}</span>
        <span
          v-if="!isAffix(tag)"
          class="close-icon w-4 h-4 inline-flex items-center justify-center rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
          @click.prevent.stop="closeSelectedTag(tag)"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </span>
      </router-link>
    </div>

    <!-- Context Menu -->
    <ul v-show="visible" :style="{left: left+'px', top: top+'px'}" class="contextmenu absolute m-0 px-1 py-1 z-50 bg-white dark:bg-gray-800 rounded-md shadow-lg border border-gray-200 dark:border-gray-700 text-sm font-medium text-gray-700 dark:text-gray-300">
      <li class="cursor-pointer px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-sm" @click="refreshSelectedTag(selectedTag)">刷新</li>
      <li v-if="!isAffix(selectedTag)" class="cursor-pointer px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-sm" @click="closeSelectedTag(selectedTag)">关闭</li>
      <li class="cursor-pointer px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-sm" @click="closeOthersTags">关闭其他</li>
      <li class="cursor-pointer px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-sm" @click="closeAllTags(selectedTag)">关闭所有</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTagsViewStore, type TagView } from '@/stores/tagsView'

const route = useRoute()
const router = useRouter()
const tagsViewStore = useTagsViewStore()

const visible = ref(false)
const top = ref(0)
const left = ref(0)
const selectedTag = ref<TagView>({})
const scrollPane = ref<HTMLElement | null>(null)

const visitedViews = computed(() => tagsViewStore.visitedViews)

const isActive = (tag: TagView) => {
  return tag.path === route.path
}

const isAffix = (tag: TagView) => {
  return tag.meta && tag.meta.affix
}

const addTags = () => {
  const { name } = route
  if (name) {
    tagsViewStore.addView(route as any)
  }
  return false
}

const initTags = () => {
  const affixTags = filterAffixTags(router.getRoutes())
  for (const tag of affixTags) {
    if (tag.name) {
      tagsViewStore.addVisitedView(tag as any)
    }
  }
}

const filterAffixTags = (routes: any[], basePath = '/') => {
  let tags: TagView[] = []
  routes.forEach(route => {
    if (route.meta && route.meta.affix) {
      const tagPath = route.path
      tags.push({
        fullPath: tagPath,
        path: tagPath,
        name: route.name,
        meta: { ...route.meta }
      })
    }
    if (route.children) {
      const tempTags = filterAffixTags(route.children, route.path)
      if (tempTags.length >= 1) {
        tags = [...tags, ...tempTags]
      }
    }
  })
  return tags
}

const refreshSelectedTag = (view: TagView) => {
  tagsViewStore.delCachedView(view).then(() => {
    const { fullPath } = view
    router.replace({ path: '/redirect' + fullPath })
  })
}

const closeSelectedTag = (view: TagView) => {
  tagsViewStore.delView(view).then(({ visitedViews }: any) => {
    if (isActive(view)) {
      toLastView(visitedViews, view)
    }
  })
}

const closeOthersTags = () => {
  if (!selectedTag.value.path) return;
  router.push(selectedTag.value)
  tagsViewStore.delOthersViews(selectedTag.value).then(() => {
    moveToCurrentTag()
  })
}

const closeAllTags = (view: TagView) => {
  tagsViewStore.delAllViews().then(({ visitedViews }: any) => {
    if (visitedViews.some((v: any) => v.path === view.path)) {
      return
    }
    toLastView(visitedViews, view)
  })
}

const toLastView = (visitedViews: TagView[], view: TagView) => {
  const latestView = visitedViews.slice(-1)[0]
  if (latestView) {
    router.push(latestView.fullPath!)
  } else {
    // defaults to dashboard
    if (view.name === 'AuthDashboard') {
      router.replace({ path: '/redirect' + view.fullPath })
    } else {
      router.push('/')
    }
  }
}

const openMenu = (tag: TagView, e: MouseEvent) => {
  const menuMinWidth = 105
  const offsetLeft = scrollPane.value?.getBoundingClientRect().left ?? 0
  const offsetWidth = scrollPane.value?.offsetWidth ?? 0
  const maxLeft = offsetWidth - menuMinWidth
  const mouseLeft = e.clientX - offsetLeft + 15

  if (mouseLeft > maxLeft) {
    left.value = maxLeft
  } else {
    left.value = mouseLeft
  }

  // Find a good general relative coordinate based on layout
  left.value = e.clientX + 10;
  top.value = e.clientY + 10;

  visible.value = true
  selectedTag.value = tag
}

const closeMenu = () => {
  visible.value = false
}

const moveToCurrentTag = () => {
  // Can implement smooth scrolling to tag if it overflows
}

watch(route, () => {
  addTags()
  moveToCurrentTag()
})

watch(visible, (value) => {
  if (value) {
    document.body.addEventListener('click', closeMenu)
  } else {
    document.body.removeEventListener('click', closeMenu)
  }
})

onMounted(() => {
  initTags()
  addTags()
})
</script>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
}
</style>
