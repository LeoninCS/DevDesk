<template>
  <div class="home">
    <h1>欢迎使用 DevDesk</h1>
    <p class="lead">请选择左侧的功能开始使用。</p>

    <section>
      <h2>功能介绍</h2>

      <h3>代码分享</h3>
      <p>填入作者、语言、保存时间，将代码粘贴后生成链接，分享给他人即可查看。</p>

      <h3>HTML 托管</h3>
      <p>粘贴你的 HTML 内容，一键生成访问链接，适合快速展示静态页面或 Demo。</p>

      <h3>Markdown 编辑</h3>
      <p>左侧编辑、右侧预览，可导出 Markdown、HTML、PDF，并通过链接协同合作。</p>

      <h3>待办事项</h3>
      <p>点击开始添加你的待办，并复制链接再次访问和修改，适合个人与团队协作。</p>

      <h3>在线 HTTP 接口测试</h3>
      <p>进入“HTTP 接口测试”功能即可使用类 Postman 的界面，支持多请求工作区、JSON 导入导出与响应头/正文查看。</p>
    </section>
  </div>
</template>

<script>
export default {
  name: "MomentWelcome",
};
</script>

<style scoped>
.home {
  max-width: 860px;
  margin: 0 auto;
  padding: 8px 10px 18px;
}

h1 {
  font-size: 26px;
  margin-bottom: 8px;
}

.lead {
  margin: 0 0 12px;
  font-size: 16px;
}

h2 {
  margin-top: 24px;
  font-size: 22px;
}

h3 {
  margin-top: 16px;
  font-size: 18px;
}

p {
  margin: 8px 0;
  line-height: 1.7;
  font-size: 16px;
}

@media (max-width: 640px) {
  .home {
    padding: 6px 6px 16px;
  }

  h1 {
    font-size: 22px;
  }

  h2 {
    font-size: 18px;
  }
}
</style>
