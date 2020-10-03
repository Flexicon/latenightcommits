<template>
  <li class="commit-entry flex items-center py-2">
    <div class="w-24 flex-shrink-0 sm:w-auto sm:flex sm:items-center">
      <a
        class="hover:opacity-75 sm:w-24 sm:flex-shrink-0"
        :href="`https://github.com/${author}`"
        rel="noreferrer noopener"
        target="_blank"
      >
        <img :src="avatar_url" :alt="`${author} avatar`" />
      </a>

      <div class="p-1 sm:px-5 text-2xs text-center">
        <div>{{ author }}</div>
        <small>{{ created_at }}</small>
      </div>
    </div>

    <div class="p-3 font-mono text-sm sm:text-base">
      <a
        class="hover:underline"
        :href="link"
        rel="noreferrer noopener"
        target="_blank"
      >
        {{ message }}
      </a>
    </div>
  </li>
</template>

<script>
export default {
  props: {
    id: String,
    author: String,
    message: String,
    created_at: String,
    avatar_url: String,
    link: String,
  },
};
</script>
