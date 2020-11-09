<template>
  <li class="commit-entry py-5 px-6 bg-white mb-4 rounded-xl">
    <div class="font-mono text-sm sm:text-base">
      <a
        class="hover:underline"
        :href="link"
        rel="noreferrer noopener"
        target="_blank"
      >
        {{ trimmedMessage }}
      </a>
    </div>

    <div class="flex flex-col-reverse sm:flex-row sm:items-center mt-3">
      <component
        :is="author ? 'a' : 'div'"
        class="flex items-center"
        :class="{ 'hover:underline': author }"
        :href="authorLink"
        :title="author"
        rel="noreferrer noopener"
        target="_blank"
      >
        <img
          class="block w-6 rounded-full"
          :src="displayImage"
          :alt="`${displayName} avatar`"
          :class="{ 'bg-gray-200': !avatar_url }"
        />

        <span class="ml-2 text-xs sm:text-sm">
          {{ displayName }}
        </span>
      </component>

      <div
        class="mb-1 sm:mb-0 sm:ml-1 text-2xs sm:text-sm text-gray-500"
        :title="timestamp"
      >
        {{ relativeTimestamp }}
      </div>
    </div>
  </li>
</template>

<script>
import { format, formatDistanceToNow } from 'date-fns/esm';

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
      return this.author || 'anonymous';
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
    relativeTimestamp() {
      return formatDistanceToNow(new Date(this.created_at), {
        addSuffix: true,
      });
    },
    timestamp() {
      return format(new Date(this.created_at), 'LLLL d, yyyy hh:mm a');
    },
  },
};
</script>
