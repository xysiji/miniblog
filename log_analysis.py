import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import json
import os

# 配置中文字体，防止图表乱码
plt.rcParams['font.sans-serif'] = ['SimHei']  # Windows用黑体，如果是Mac请改为 'Arial Unicode MS'
plt.rcParams['axes.unicode_minus'] = False

# 【核心修改】：精准匹配你的实际日志文件名
LOG_FILE_NAME = 'data_analysis.log'

def load_data(filepath=LOG_FILE_NAME):
    print(f"🔍 正在当前目录读取分布式存储访问日志: {os.path.abspath(filepath)}")
    data = []
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            for line in f:
                if line.strip():
                    try:
                        data.append(json.loads(line.strip()))
                    except json.JSONDecodeError:
                        continue # 跳过无法解析的脏数据行
    except FileNotFoundError:
        print(f"❌ 未找到日志文件 {filepath}，请确认后端项目已产生埋点日志！")
        return pd.DataFrame()
        
    df = pd.DataFrame(data)
    
    # 【新增的数据防呆与清洗机制】
    if df.empty:
        print("❌ 日志文件为空！")
        return df
        
    if 'post_id' not in df.columns:
        print("❌ 致命错误：当前日志中缺少核心业务字段 'post_id'！")
        print(f"⚠️ 当前读取到的字段有：{list(df.columns)}")
        print("👉 解决方案：请先删除旧的 data_analysis.log 文件，然后去网页上进行点赞和评论操作，生成全新的结构化埋点数据！")
        return pd.DataFrame()
        
    # 清洗掉可能混入的没有 post_id 的旧数据行
    df = df.dropna(subset=['post_id'])
    
    print(f"✅ 成功加载并清洗日志数据，共计 {len(df)} 条有效埋点记录。")
    return df

def generate_visualizations():
    df = load_data()
    if df.empty:
        return

    # 1. 转换时间格式 (为了后续的时间序列分析)
    df['timestamp'] = pd.to_datetime(df['timestamp'])
    
    # 初始化画布 (16:10 宽屏比例，适合放在毕业论文中)
    fig = plt.figure(figsize=(16, 10))
    fig.suptitle('微型博客分布式系统 - 存储访问与互动行为分析图谱', fontsize=20, fontweight='bold')

    # 图表 1：热门博文 Top 10 互动量 (模拟热点数据探测，点题分布式存储的数据倾斜)
    ax1 = fig.add_subplot(221)
    top_posts = df['post_id'].astype(str).value_counts().head(10)
    sns.barplot(x=top_posts.values, y=top_posts.index, palette='viridis', ax=ax1)
    ax1.set_title('Top 10 热门博文 (数据倾斜/热点探测)')
    ax1.set_xlabel('互动总次数 (点赞+评论)')
    ax1.set_ylabel('Post ID')

    # 图表 2：用户行为漏斗 / 动作类型分布 (分析读写比例)
    ax2 = fig.add_subplot(222)
    action_counts = df['action_type'].value_counts()
    ax2.pie(action_counts, labels=action_counts.index, autopct='%1.1f%%', colors=['#ff9999','#66b3ff'], startangle=90, shadow=True)
    ax2.set_title('用户交互行为类型分布 (读写比例参考)')

    # 图表 3：接口响应耗时分布 (验证微服务与缓存架构的性能)
    ax3 = fig.add_subplot(223)
    sns.histplot(df['cost_time_ms'], bins=20, kde=True, color='teal', ax=ax3)
    ax3.set_title('分布式架构 API 接口响应耗时分布 (ms)')
    ax3.set_xlabel('响应时间 (毫秒)')
    ax3.set_ylabel('频次')
    
    # 图表 4：流量时间序列图 (削峰填谷、并发监控)
    ax4 = fig.add_subplot(224)
    # 按分钟重采样，如果你的测试数据是在同一分钟内产生的，可以用 'S' 按秒重采样
    time_series = df.set_index('timestamp').resample('1T').size() 
    time_series.plot(kind='line', marker='o', color='coral', ax=ax4)
    ax4.set_title('系统并发访问时间序列 (流量峰值监控)')
    ax4.set_xlabel('时间')
    ax4.set_ylabel('请求次数 (QPM)')

    # 调整排版并保存高清大图
    plt.tight_layout(rect=[0, 0.03, 1, 0.95])
    output_filename = 'architecture_data_analysis.png'
    plt.savefig(output_filename, dpi=300)
    print(f"🎉 深度分析完成！图表已保存为高清图片：{output_filename}，快去把它贴进你的毕业论文里吧！")
    
    # 在运行窗口中展示
    plt.show()

if __name__ == "__main__":
    generate_visualizations()