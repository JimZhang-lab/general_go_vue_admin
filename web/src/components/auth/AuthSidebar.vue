<template>
  <aside
    :class="[
      'fixed mt-16 flex flex-col lg:mt-0 top-0 px-5 left-0 bg-white dark:bg-gray-900 dark:border-gray-800 text-gray-900 h-screen transition-all duration-300 ease-in-out z-99999 border-r border-gray-200',
      {
        'lg:w-[290px]': isExpanded || isMobileOpen || isHovered,
        'lg:w-[90px]': !isExpanded && !isHovered,
        'translate-x-0 w-[290px]': isMobileOpen,
        '-translate-x-full': !isMobileOpen,
        'lg:translate-x-0': true,
      },
    ]"
    @mouseenter="!isExpanded && (isHovered = true)"
    @mouseleave="isHovered = false"
  >
    <!-- Logo Section -->
    <div
      :class="[
        'py-8 flex',
        !isExpanded && !isHovered ? 'lg:justify-center' : 'justify-start',
      ]"
    >
      <router-link to="/auth/dashboard">
        <img
          v-if="isExpanded || isHovered || isMobileOpen"
          class="dark:hidden"
          src="/images/logo/logo.svg"
          alt="Logo"
          width="150"
          height="40"
        />
        <img
          v-if="isExpanded || isHovered || isMobileOpen"
          class="hidden dark:block"
          src="/images/logo/logo-dark.svg"
          alt="Logo"
          width="150"
          height="40"
        />
        <img
          v-else
          src="/images/logo/logo-icon.svg"
          alt="Logo"
          width="32"
          height="32"
        />
      </router-link>
    </div>

    <!-- Navigation Menu -->
    <div class="flex flex-col overflow-y-auto duration-300 ease-linear no-scrollbar">
      <nav class="mb-6">
        <div class="flex flex-col" style="gap: 1rem;">
          <div v-for="(menuGroup, groupIndex) in authMenuGroups" :key="groupIndex">
            <h2
              :class="[
                'mb-4 text-xs uppercase flex leading-[20px] text-gray-400',
                !isExpanded && !isHovered
                  ? 'lg:justify-center'
                  : 'justify-start',
              ]"
            >
              <template v-if="isExpanded || isHovered || isMobileOpen">
                {{ menuGroup.title }}
              </template>
              <HorizontalDots v-else />
            </h2>
            <ul class="flex flex-col" style="gap: 1rem;">
              <li v-for="(item, index) in menuGroup.items" :key="item.name">
                <!-- 有子菜单的项目 -->
                <button
                  v-if="item.subItems"
                  @click="toggleSubmenu(groupIndex, index)"
                  :class="[
                    'menu-item group w-full',
                    {
                      'menu-item-active': isSubmenuOpen(groupIndex, index),
                      'menu-item-inactive': !isSubmenuOpen(groupIndex, index),
                    },
                    !isExpanded && !isHovered
                      ? 'lg:justify-center'
                      : 'lg:justify-start',
                  ]"
                >
                  <span
                    :class="[
                      isSubmenuOpen(groupIndex, index)
                        ? 'menu-item-icon-active'
                        : 'menu-item-icon-inactive',
                    ]"
                  >
                    <component :is="item.icon" />
                  </span>
                  <span
                    v-if="isExpanded || isHovered || isMobileOpen"
                    class="menu-item-text"
                    >{{ item.name }}</span
                  >
                  <ChevronDownIcon
                    v-if="isExpanded || isHovered || isMobileOpen"
                    :class="[
                      'ml-auto w-5 h-5 transition-transform duration-200',
                      {
                        'rotate-180 text-brand-500': isSubmenuOpen(
                          groupIndex,
                          index
                        ),
                      },
                    ]"
                  />
                </button>
                <!-- 直接链接的项目 -->
                <router-link
                  v-else-if="item.path"
                  :to="item.path"
                  :class="[
                    'menu-item group',
                    {
                      'menu-item-active': isActive(item.path),
                      'menu-item-inactive': !isActive(item.path),
                    },
                  ]"
                >
                  <span
                    :class="[
                      isActive(item.path)
                        ? 'menu-item-icon-active'
                        : 'menu-item-icon-inactive',
                    ]"
                  >
                    <component :is="item.icon" />
                  </span>
                  <span
                    v-if="isExpanded || isHovered || isMobileOpen"
                    class="menu-item-text"
                    >{{ item.name }}</span
                  >
                </router-link>
                <!-- 子菜单展开 -->
                <transition
                  @enter="startTransition"
                  @after-enter="endTransition"
                  @before-leave="startTransition"
                  @after-leave="endTransition"
                >
                  <div
                    v-show="
                      isSubmenuOpen(groupIndex, index) &&
                      (isExpanded || isHovered || isMobileOpen)
                    "
                  >
                    <ul class="mt-2 space-y-1 ml-9">
                      <li v-for="subItem in item.subItems" :key="subItem.name">
                        <router-link
                          :to="subItem.path"
                          :class="[
                            'submenu-item group',
                            {
                              'submenu-item-active': isActive(subItem.path),
                              'submenu-item-inactive': !isActive(subItem.path),
                            },
                          ]"
                        >
                          <span class="submenu-item-text">{{ subItem.name }}</span>
                        </router-link>
                      </li>
                    </ul>
                  </div>
                </transition>
              </li>
            </ul>
          </div>
        </div>
      </nav>

      <!-- Auth Widget -->
      <AuthSidebarWidget v-if="isExpanded || isHovered || isMobileOpen" />
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import { useSidebar } from "@/composables/useSidebar";
import {
  ChevronDownIcon,
  HorizontalDots,
  GridIcon,
  UserCircleIcon,
  SettingsIcon,
  UserGroupIcon,
  PlugInIcon,
  DocsIcon,
  SettingsIcon as CogIcon,
  BarChartIcon,
} from "@/icons";
import AuthSidebarWidget from "./AuthSidebarWidget.vue";

const route = useRoute();
const { isExpanded, isMobileOpen, isHovered, openSubmenu } = useSidebar();

// Auth模块专用菜单配置
const authMenuGroups = [
  {
    title: "权限管理",
    items: [
      {
        icon: BarChartIcon,
        name: "权限总览",
        path: "/auth/dashboard",
      },
      {
        icon: UserGroupIcon,
        name: "管理员管理",
        path: "/auth/admin",
      },
      {
        icon: SettingsIcon,
        name: "角色管理",
        path: "/auth/role",
      },
      {
        icon: PlugInIcon,
        name: "权限管理",
        path: "/auth/permission",
      },
      {
        icon: DocsIcon,
        name: "系统日志",
        path: "/auth/logs",
      },
    ],
  },
  {
    title: "个人中心",
    items: [
      {
        icon: UserCircleIcon,
        name: "个人资料",
        path: "/auth/profile",
      },
      {
        icon: CogIcon,
        name: "系统设置",
        subItems: [
          { name: "基础设置", path: "/auth/settings/basic" },
          { name: "安全设置", path: "/auth/settings/security" },
          { name: "通知设置", path: "/auth/settings/notification" },
        ],
      },
    ],
  },
  {
    title: "其他功能",
    items: [
      {
        icon: GridIcon,
        name: "返回主页",
        path: "/",
      },
    ],
  },
];

const isActive = (path) => route.path === path;

const toggleSubmenu = (groupIndex, itemIndex) => {
  const key = `${groupIndex}-${itemIndex}`;
  openSubmenu.value = openSubmenu.value === key ? null : key;
};

const isAnySubmenuRouteActive = computed(() => {
  return authMenuGroups.some((group) =>
    group.items.some(
      (item) =>
        item.subItems && item.subItems.some((subItem) => isActive(subItem.path))
    )
  );
});

const isSubmenuOpen = (groupIndex, itemIndex) => {
  const key = `${groupIndex}-${itemIndex}`;
  return (
    openSubmenu.value === key ||
    (isAnySubmenuRouteActive.value &&
      authMenuGroups[groupIndex].items[itemIndex].subItems?.some((subItem) =>
        isActive(subItem.path)
      ))
  );
};

const startTransition = (el) => {
  el.style.height = "auto";
  const height = el.scrollHeight;
  el.style.height = "0px";
  el.offsetHeight; // force reflow
  el.style.height = height + "px";
};

const endTransition = (el) => {
  el.style.height = "";
};
</script>

<style scoped>
/* 菜单项样式 */
.menu-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.625rem 0.75rem;
  border-radius: 0.5rem;
  transition: all 0.2s;
  font-weight: 500;
  font-size: 0.875rem;
}

.menu-item-active {
  background-color: rgb(239 246 255);
  color: rgb(29 78 216);
}

.menu-item-inactive {
  color: rgb(75 85 99);
}

.menu-item-inactive:hover {
  background-color: rgb(243 244 246);
  color: rgb(17 24 39);
}

.menu-item-icon-active {
  color: rgb(37 99 235);
}

.menu-item-icon-inactive {
  color: rgb(107 114 128);
}

.menu-item-text {
  font-weight: 500;
}

/* 子菜单项样式 */
.submenu-item {
  display: flex;
  align-items: center;
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  transition: all 0.2s;
  font-size: 0.875rem;
}

.submenu-item-active {
  background-color: rgb(239 246 255);
  color: rgb(29 78 216);
}

.submenu-item-inactive {
  color: rgb(75 85 99);
}

.submenu-item-inactive:hover {
  background-color: rgb(243 244 246);
  color: rgb(17 24 39);
}

.submenu-item-text {
  font-weight: 500;
}

/* 暗色模式 */
@media (prefers-color-scheme: dark) {
  .menu-item-active {
    background-color: rgba(59, 130, 246, 0.15);
    color: rgb(96 165 250);
  }

  .menu-item-inactive {
    color: rgb(156 163 175);
  }

  .menu-item-inactive:hover {
    background-color: rgb(31 41 55);
    color: rgb(209 213 219);
  }

  .menu-item-icon-active {
    color: rgb(96 165 250);
  }

  .menu-item-icon-inactive {
    color: rgb(156 163 175);
  }

  .submenu-item-active {
    background-color: rgba(59, 130, 246, 0.15);
    color: rgb(96 165 250);
  }

  .submenu-item-inactive {
    color: rgb(156 163 175);
  }

  .submenu-item-inactive:hover {
    background-color: rgb(31 41 55);
    color: rgb(209 213 219);
  }
}
</style>
