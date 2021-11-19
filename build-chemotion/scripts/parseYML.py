import click
import yaml
import os.path


@click.group()
def cli():
    pass


@click.command()
@click.option("--prefix", default="",
              help="string prefix for imported strings")
@click.option("--upper", is_flag=True, default=False,
              help="returned keywords ""should be uppercase letters")
@click.option("--content", is_flag=True, default=False,
              help="returned keywords should only contain the returned value")
@click.option("--title", is_flag=True, default=False,
              help="returned keywords should only contain the titles")
@click.option("--collect", is_flag=True, default=False,
              help="return all values")
@click.argument("file", nargs=1)
@click.argument("keywords", nargs=1)
def read(prefix, upper, content, title, collect, file, keywords):
    """Read KEYWORDS (in UPPER case) from a FILE and add the given PREFIX"""
    if os.path.isfile(file):
        stream = open(file, "r")
        fileContent = yaml.safe_load(stream)
    
        if collect:
            if len(keywords.split(".")) == 2:
                for v in fileContent[keywords.split(".")[0]].values():
                    click.echo(v[keywords.split(".")[1]])
        else:
            notFound = False
            currentKeyword = ""
            for keyword in keywords.split("."):
                if type(fileContent) is dict:
                    currentKeyword = keyword
                    if fileContent.get(keyword):
                        fileContent = fileContent[keyword]
                    else:
                        notFound = True
                else:
                    notFound = True
    
            if not notFound:
                if type(fileContent) is dict:
                    for item in fileContent:
                        if type(item) is str:
                            if upper:
                                click.echo(prefix + item.upper() +
                                           "=" + str(fileContent[item]))
                            elif title:
                                click.echo(item)
                            elif content:
                                click.echo(str(fileContent[item]))
                            else:
                                click.echo(prefix + item + "=" +
                                           str(fileContent[item]))
    
                else:
                    if upper:
                        click.echo(prefix + currentKeyword.upper() +
                                   "=" + str(fileContent))
                    elif title:
                        click.echo(currentKeyword)
                    elif content:
                        click.echo(str(fileContent))
                    else:
                        click.echo(prefix + currentKeyword +
                                   "=" + str(fileContent))


cli.add_command(read)

if __name__ == '__main__':
    cli()
