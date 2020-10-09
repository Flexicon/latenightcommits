<template>
  <li class="commit-entry flex items-center py-2">
    <div class="w-24 flex-shrink-0 sm:w-auto sm:flex sm:items-center">
      <component
        :is="author ? 'a' : 'div'"
        class="hover:opacity-75 sm:w-24 sm:flex-shrink-0"
        :href="`https://github.com/${author}`"
        rel="noreferrer noopener"
        target="_blank"
      >
        <img :src="displayImage" :alt="`${displayName} avatar`" :class="{ 'bg-gray-100': !avatar_url }" />
      </component>

      <div class="p-1 sm:px-5 text-2xs text-center">
        <div>{{ displayName }}</div>
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
        {{ trimmedMessage }}
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
  computed: {
    displayName() {
      return this.author || '[REDACTED]';
    },
    displayImage() {
      return (
        this.avatar_url ||
        require('@/assets/redacted_user.svg')
      );
    },
    trimmedMessage() {
      const maxLength = 120;

      return this.message.length > maxLength
        ? `${this.message.slice(0, maxLength)}...`
        : this.message;
    },
  },
};
</script>
