import json
import pandas as pd
import matplotlib.pyplot as plt
import matplotlib
import os

# 设置中文字体，防止图表中的中文显示为方块
matplotlib.rcParams['font.sans-serif'] = ['SimHei'] # Windows系统使用黑体
matplotlib.rcParams['axes.unicode_minus'] = False

def analyze_logs(log_file):
    print(f"🔍 正在当前目录查找日志文件: {os.path.abspath(log_file)}")
    
    if not os.path.exists(log_file):
        print("❌ 错误：未找到日志文件！请确认是否在正确目录下运行此脚本。")
        return

    data = []
    try:
        with open(log_file, 'r', encoding='utf-8') as f:
            for line in f:
                if line.strip():
                    data.append(json.loads(line.strip()))
    except Exception as e:
        print(f"❌ 读取日志出错: {e}")
        return

    df = pd.DataFrame(data)
    if df.empty:
        print("❌ 日志文件为空！")
        return

    # 将微秒(duration)转换为毫秒(ms)
    df['duration_ms'] = df['duration'] / 1000.0

    print("\n✅ --- 成功读取原始数据 (最新5条) ---")
    print(df.tail())

    # 按接口路径分组统计：平均耗时和请求次数
    stats = df.groupby('path').agg(
        avg_duration=('duration_ms', 'mean'),
        request_count=('path', 'count')
    ).reset_index()

    print("\n✅ --- 接口性能统计摘要 ---")
    print(stats)

    # 生成可视化图表
    plt.figure(figsize=(10, 6))
    # 动态分配柱子颜色
    colors = ['#4CAF50', '#2196F3', '#FFC107', '#E91E63'][:len(stats)]
    bars = plt.bar(stats['path'], stats['avg_duration'], color=colors)
    
    plt.title('微型博客核心接口平均响应耗时分析 (多级缓存性能验证)', fontsize=14)
    plt.xlabel('API 接口路径', fontsize=12)
    plt.ylabel('平均响应耗时 (毫秒)', fontsize=12)
    plt.grid(axis='y', linestyle='--', alpha=0.7)

    # 在柱子上添加具体数值标签 (包含耗时和请求次数)
    for bar in bars:
        yval = bar.get_height()
        req_count = int(stats.loc[stats['avg_duration'] == yval, 'request_count'].values[0])
        label_text = f"{yval:.2f} ms\n({req_count}次请求)"
        plt.text(bar.get_x() + bar.get_width()/2, yval + 0.1, label_text, ha='center', va='bottom', fontsize=10)

    # 获取绝对路径并保存图片 (去掉了 plt.show 防止环境卡死)
    output_img = os.path.abspath("api_performance_chart.png")
    plt.savefig(output_img, dpi=300, bbox_inches='tight')
    
    print(f"\n🎉 成功！可视化图表已生成，请复制下方路径到文件管理器中打开：")
    print(f"👉 {output_img}")

if __name__ == "__main__":
    analyze_logs("data_analysis.log")