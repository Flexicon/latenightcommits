<template>
  <li
    class="commit-entry flex items-center p-3 pr-0 bg-white"
  >
    <div class="w-24 flex-shrink-0">
      <component
        :is="author ? 'a' : 'div'"
        class="block"
        :class="{ 'hover:opacity-75': author }"
        :href="authorLink"
        rel="noreferrer noopener"
        target="_blank"
      >
        <img
          :src="displayImage"
          :alt="`${displayName} avatar`"
          :class="{ 'bg-gray-200': !avatar_url }"
        />
      </component>

      <div class="p-1 text-2xs text-center">
        <component
          :is="author ? 'a' : 'div'"
          class="block w-23 truncate"
          :class="{ 'hover:underline': author }"
          :href="authorLink"
          :title="author"
          rel="noreferrer noopener"
          target="_blank"
        >
          {{ displayName }}
        </component>
        <small>{{ created_at }}</small>
      </div>
    </div>

    <div class="p-4 sm:px-8 font-mono text-sm sm:text-base">
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
      return this.avatar_url || require('@/assets/redacted_user.svg');
    },
    authorLink() {
      return `https://github.com/${this.author}`;
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
