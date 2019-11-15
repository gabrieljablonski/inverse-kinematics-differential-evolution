import matplotlib.pyplot as plt
from matplotlib.widgets import Slider, TextBox
from mpl_toolkits.mplot3d import Axes3D


def main(filename):
    with open(filename) as f:
        points = f.readlines()

    target = list(map(float, points[0].split(',')))
    generations = []
    for line in points[1:]:
        links = line.split('\t')
        links_as_points = []
        for link in links:
            links_as_points.append(list(map(float, link.split(','))))
        generations.append(links_as_points)

    fig = plt.figure()
    ax = fig.add_subplot(111, projection='3d')
    ax.grid(False)
    ax.set_xlim((-.3, .3))
    ax.set_ylim((-.3, .3))
    ax.set_zlim((-.3, .3))

    gen_slider_ax = fig.add_axes([0.2, .95, .13, .03], label='Generation')
    gen_slider = Slider(gen_slider_ax,
                        'Generation',
                        1, len(generations),
                        valfmt='%d', valinit=1, valstep=1)

    fitness_ax = fig.add_axes([0.2, .88, .13, .03], label='Fitness')
    fitness_text = TextBox(fitness_ax, '', str('0.00000'))

    ax.plot(*zip(target), 'rx', markersize=10)
    colors = 'black', 'red', 'yellow', 'blue'
    link_plots = [
        ax.plot((), (), (), color=color, linewidth=4, solid_capstyle='round')[0]
        for color in colors
    ]
    gen = 0

    def update_plot(_):
        generation = generations[int(gen_slider.val)-1]
        link0, *links1_4 = generation
        last = link0
        for i, link in enumerate(links1_4):
            x_data, y_data, z_data = zip(last, link)
            link_plots[i].set_xdata(x_data)
            link_plots[i].set_ydata(y_data)
            link_plots[i].set_3d_properties(z_data)
            fig.canvas.draw_idle()
            last = link

        def fitness(point, target):
            x, y, z = point
            tx, ty, tz = target
            return ((x-tx)**2+(y-ty)**2+(z-tz)**2)**.5
        fitness_text.set_val(f'{fitness(last, target):.5f}')

    gen_slider.on_changed(update_plot)

    plt.show()


if __name__ == '__main__':
    import sys

    file_name = sys.argv[1]
    main(file_name)
