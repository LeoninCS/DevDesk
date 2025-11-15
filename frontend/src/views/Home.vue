<template>
  <section>
    <h1 style="margin:0 0 12px 0;">功能总览</h1>
    <p class="helper">点击下面的功能卡片进入对应模块。当前已开放 <strong>代码分享</strong>。</p>
    <div style="height:14px"></div>
    <div class="grid">
      <FeatureCard title="代码分享" description="上传代码，生成唯一哈希链接，分享只读页面。" badge="可用" to="/codeshare" />
      <FeatureCard title="加密/解密文本" description="上传文本和密钥，可选择加密或者解密文本。" badge="可用" to="/textcrypto" />
      <FeatureCard title="文本短链（预留）" description="生成一次性文本访问链接。" badge="开发中">
        <template #action>
          <button class="btn" disabled>敬请期待</button>
        </template>
      </FeatureCard>
      <FeatureCard title="文件分享（预留）" description="临时文件分享与过期销毁。" badge="开发中">
        <template #action>
          <button class="btn" disabled>敬请期待</button>
        </template>
      </FeatureCard>
    </div>
  </section>
</template>

<script setup>
import FeatureCard from '../components/FeatureCard.vue'
</script>
