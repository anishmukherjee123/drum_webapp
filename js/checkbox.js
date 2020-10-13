$(document).ready(function() {
  $('.input_class_checkbox').each(function(){
    $(this).hide().after('<div class="class_checkbox" />');

  });

  $('.class_checkbox').on('click',function(){
    $(this).toggleClass('checked').prev().prop('checked',$(this).is('.checked'))
  });
});