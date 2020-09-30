function handleTotal(total, get){

  $("#total").removeClass('red green yellow')

    switch (true) {
      case (total > 2000):
        $("#total").addClass('red')
        break;
    case (total > 1800):
        $("#total").addClass('yellow')
        break;
      default:
        $("#total").addClass('green')
    }

    $("#total").html(total)
}
