class AIModel {
  late String name;
  late String provider;

  AIModel({required this.name, required this.provider});

  AIModel.fromJson(Map<String, dynamic> json) {
    name = json['Name'];
    provider = json['Provider'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['Name'] = name;
    data['Provider'] = provider;
    return data;
  }
}

class AIMCPTool {
  late String name;
  late String description;

  AIMCPTool({required this.name, required this.description});

  AIMCPTool.fromJson(Map<String, dynamic> json) {
    name = json['Name'] ?? "";
    description = json['Description'] ?? "";
  }
}
